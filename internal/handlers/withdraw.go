package handlers

import (
	"finance/database"
	"finance/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func HandleWithdrawal(c fiber.Ctx) error {

	var data map[string]string

	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	var totalAmount float64
	var user models.User

	database.DB.First(&user, "user_id = ?", data["userId"])

	database.DB.Table("deposits").Where("user_Id = ?", data["userId"]).Group("user_Id").Pluck("SUM(amount)", &totalAmount)

	amountString := data["amount"]

	amount, _ := strconv.ParseFloat(amountString, 64)

	withdraw := models.Withdraw{
		UserID:      data["userId"],
		Name:        data["name"],
		Amount:      amount,
		Currency:    data["currencyy"],
		TotalAmount: totalAmount - amount,
	}

	//update deposits table
	deposit := models.Deposit{
		UserID:      data["userId"],
		Name:        data["name"],
		Currency:    data["currencyy"],
		Withdrawal:  -amount,
		TotalAmount: totalAmount - amount,
	}

	database.DB.Create(&withdraw)
	database.DB.Create(&deposit)
	return c.JSON(withdraw)
}
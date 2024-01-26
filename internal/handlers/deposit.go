package handlers

import (
	"finance/database"
	"finance/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func NewDeposit(userId string, name string, amount float64, currency string) *models.Deposit {
	return &models.Deposit{
		UserID:   userId,
		Name:     name,
		Amount:   amount,
		Currency: currency,
	}
}

func HandleCreateDeposit(c fiber.Ctx) error {

	var data map[string]string

	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	var totalAmount float64

	database.DB.Table("deposits").Where("user_Id = ?", data["userId"]).Group("user_Id").Pluck("SUM(amount)", &totalAmount)
	log.Println("total amount", totalAmount)

	amountString := data["amount"]

	amount, _ := strconv.ParseFloat(amountString, 64)

	deposit := models.Deposit{
		UserID:   data["userId"],
		Name:     data["name"],
		Amount:   amount + totalAmount,
		Currency: data["currency"],
	}

	log.Println("here is the Deposit", deposit)
	newDeposit := NewDeposit(deposit.UserID, deposit.Name, deposit.Amount, deposit.Currency)

	database.DB.Create(newDeposit)
	return c.JSON(newDeposit)
}
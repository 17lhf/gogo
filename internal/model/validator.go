package model

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidators registers custom validation rules for enum types.
func RegisterValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	v.RegisterValidation("userstatus", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(UserStatus)
		if !ok {
			return false
		}
		switch val {
		case UserStatusEnabled, UserStatusDisabled, UserStatusLocked:
			return true
		}
		return false
	})

	v.RegisterValidation("terminalstatus", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(TerminalStatus)
		if !ok {
			return false
		}
		switch val {
		case TerminalStatusOffline, TerminalStatusOnline, TerminalStatusDisabled, TerminalStatusEnabled:
			return true
		}
		return false
	})

	v.RegisterValidation("menutype", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(MenuType)
		if !ok {
			return false
		}
		switch val {
		case MenuTypeDirectory, MenuTypePage, MenuTypeButton:
			return true
		}
		return false
	})

	v.RegisterValidation("logstatus", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(LogStatus)
		if !ok {
			return false
		}
		switch val {
		case LogStatusSuccess, LogStatusFailure:
			return true
		}
		return false
	})
}

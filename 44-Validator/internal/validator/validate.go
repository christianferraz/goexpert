package validator

import (
	"context"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	EmailRX   = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PhoneRX   = regexp.MustCompile(`^\(?\d{2}\)?\s?\d{4,5}-?\d{4}$`)
	cpfDigits = regexp.MustCompile(`[^0-9]`)
)

type Validator interface {
	Valid(context.Context) Evaluator
}

type Evaluator map[string]string

func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(Evaluator)
	}
	if _, ok := (*e)[key]; !ok {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(ok bool, key, message string) {
	if !ok {
		e.AddFieldError(key, message)
	}
}

// StructValidator é o validador principal que usa struct tags
type StructValidator struct {
	data any
}

// NewStructValidator cria um novo validador para struct
func NewStructValidator(data any) *StructValidator {
	return &StructValidator{data: data}
}

// Valid implementa a interface Validator
func (sv *StructValidator) Valid(ctx context.Context) Evaluator {
	evaluator := make(Evaluator)
	sv.validateStruct(ctx, &evaluator)
	return evaluator
}

// validateStruct processa as struct tags e valida os campos
func (sv *StructValidator) validateStruct(ctx context.Context, evaluator *Evaluator) {
	v := reflect.ValueOf(sv.data)
	t := reflect.TypeOf(sv.data)

	// Se for ponteiro, pega o elemento
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	// Só processa structs
	if v.Kind() != reflect.Struct {
		return
	}

	// Itera pelos campos da struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name

		// Pula campos não exportados
		if !field.CanInterface() {
			continue
		}

		// Obtém a tag de validação
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" || validateTag == "-" {
			continue
		}

		// Processa as regras da tag
		sv.processValidationRules(ctx, evaluator, fieldName, field, validateTag)
	}
}

// processValidationRules processa as regras de validação de um campo
func (sv *StructValidator) processValidationRules(ctx context.Context, evaluator *Evaluator, fieldName string, field reflect.Value, rules string) {
	// Divide as regras por vírgula
	rulesList := strings.Split(rules, ",")

	// Verifica se tem omitempty e se o campo está vazio
	hasOmitEmpty := false
	for _, rule := range rulesList {
		if strings.TrimSpace(rule) == "omitempty" {
			hasOmitEmpty = true
			break
		}
	}

	// Se tem omitempty e o campo está vazio, pula todas as validações
	if hasOmitEmpty && sv.isEmpty(field) {
		return
	}

	for _, rule := range rulesList {
		rule = strings.TrimSpace(rule)
		if rule == "" || rule == "omitempty" {
			continue // Pula regra vazia ou omitempty (já processada acima)
		}

		// Separa a regra do parâmetro (ex: "min=2")
		parts := strings.SplitN(rule, "=", 2)
		ruleName := parts[0]
		var param string
		if len(parts) > 1 {
			param = parts[1]
		}

		sv.applyValidationRule(evaluator, fieldName, field, ruleName, param)
	}
}

// applyValidationRule aplica uma regra específica de validação
func (sv *StructValidator) applyValidationRule(evaluator *Evaluator, fieldName string, field reflect.Value, ruleName, param string) {
	switch ruleName {
	case "required":
		sv.validateRequired(evaluator, fieldName, field)
	case "min":
		if param != "" {
			if min, err := strconv.Atoi(param); err == nil {
				sv.validateMin(evaluator, fieldName, field, min)
			}
		}
	case "max":
		if param != "" {
			if max, err := strconv.Atoi(param); err == nil {
				sv.validateMax(evaluator, fieldName, field, max)
			}
		}
	case "len":
		if param != "" {
			if length, err := strconv.Atoi(param); err == nil {
				sv.validateLen(evaluator, fieldName, field, length)
			}
		}
	case "email":
		sv.validateEmail(evaluator, fieldName, field)
	case "cpf":
		sv.validateCPF(evaluator, fieldName, field)
	case "phone":
		sv.validatePhone(evaluator, fieldName, field)
	case "numeric":
		sv.validateNumeric(evaluator, fieldName, field)
	case "alpha":
		sv.validateAlpha(evaluator, fieldName, field)
	case "alphanum":
		sv.validateAlphaNum(evaluator, fieldName, field)
	case "url":
		sv.validateURL(evaluator, fieldName, field)
		// Removemos o case "omitempty" daqui pois é tratado antes
	}
}

// Funções de validação específicas
func (sv *StructValidator) validateRequired(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if sv.isEmpty(field) {
		evaluator.AddFieldError(fieldName, fieldName+" é obrigatório")
	}
}

func (sv *StructValidator) validateMin(evaluator *Evaluator, fieldName string, field reflect.Value, min int) {
	switch field.Kind() {
	case reflect.String:
		if utf8.RuneCountInString(field.String()) < min {
			evaluator.AddFieldError(fieldName, fieldName+" deve ter pelo menos "+strconv.Itoa(min)+" caracteres")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(min) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser pelo menos "+strconv.Itoa(min))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() < uint64(min) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser pelo menos "+strconv.Itoa(min))
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() < float64(min) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser pelo menos "+strconv.Itoa(min))
		}
	}
}

func (sv *StructValidator) validateMax(evaluator *Evaluator, fieldName string, field reflect.Value, max int) {
	switch field.Kind() {
	case reflect.String:
		if utf8.RuneCountInString(field.String()) > max {
			evaluator.AddFieldError(fieldName, fieldName+" deve ter no máximo "+strconv.Itoa(max)+" caracteres")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(max) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser no máximo "+strconv.Itoa(max))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() > uint64(max) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser no máximo "+strconv.Itoa(max))
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() > float64(max) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser no máximo "+strconv.Itoa(max))
		}
	}
}

func (sv *StructValidator) validateLen(evaluator *Evaluator, fieldName string, field reflect.Value, length int) {
	switch field.Kind() {
	case reflect.String:
		if utf8.RuneCountInString(field.String()) != length {
			evaluator.AddFieldError(fieldName, fieldName+" deve ter exatamente "+strconv.Itoa(length)+" caracteres")
		}
	}
}

func (sv *StructValidator) validateEmail(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		if !ValidEmail(field.String()) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ter formato de email válido")
		}
	}
}

func (sv *StructValidator) validateCPF(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		if !ValidCPF(field.String()) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser um CPF válido")
		}
	}
}

func (sv *StructValidator) validatePhone(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		if !ValidPhone(field.String()) {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser um telefone válido")
		}
	}
}

func (sv *StructValidator) validateNumeric(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		str := field.String()
		for _, char := range str {
			if char < '0' || char > '9' {
				evaluator.AddFieldError(fieldName, fieldName+" deve conter apenas números")
				return
			}
		}
	}
}

func (sv *StructValidator) validateAlpha(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		str := field.String()
		for _, char := range str {
			if !isAlpha(char) {
				evaluator.AddFieldError(fieldName, fieldName+" deve conter apenas letras")
				return
			}
		}
	}
}

// isAlpha verifica se o caractere é uma letra (incluindo acentos e ç)
func isAlpha(char rune) bool {
	// Letras básicas (a-z, A-Z) e espaço
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == ' ' {
		return true
	}

	// Caracteres acentuados comuns em português
	switch char {
	// Minúsculas com acento
	case 'á', 'à', 'â', 'ã', 'ä':
		return true
	case 'é', 'è', 'ê', 'ë':
		return true
	case 'í', 'ì', 'î', 'ï':
		return true
	case 'ó', 'ò', 'ô', 'õ', 'ö':
		return true
	case 'ú', 'ù', 'û', 'ü':
		return true
	case 'ç':
		return true
	case 'ñ':
		return true

	// Maiúsculas com acento
	case 'Á', 'À', 'Â', 'Ã', 'Ä':
		return true
	case 'É', 'È', 'Ê', 'Ë':
		return true
	case 'Í', 'Ì', 'Î', 'Ï':
		return true
	case 'Ó', 'Ò', 'Ô', 'Õ', 'Ö':
		return true
	case 'Ú', 'Ù', 'Û', 'Ü':
		return true
	case 'Ç':
		return true
	case 'Ñ':
		return true
	}

	return false
}

func (sv *StructValidator) validateAlphaNum(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		str := field.String()
		for _, char := range str {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
				evaluator.AddFieldError(fieldName, fieldName+" deve conter apenas letras e números")
				return
			}
		}
	}
}

func (sv *StructValidator) validateURL(evaluator *Evaluator, fieldName string, field reflect.Value) {
	if field.Kind() == reflect.String {
		url := field.String()
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			evaluator.AddFieldError(fieldName, fieldName+" deve ser uma URL válida (http:// ou https://)")
		}
	}
}

// isEmpty verifica se um campo está vazio
func (sv *StructValidator) isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return strings.TrimSpace(field.String()) == ""
	case reflect.Bool:
		return !field.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Pointer, reflect.Interface:
		return field.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return field.Len() == 0
	default:
		return field.IsZero()
	}
}

// Funções auxiliares já existentes
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func ValidCPF(cpf string) bool {
	// Remove tudo que não é número
	cpf = cpfDigits.ReplaceAllString(cpf, "")

	// Agora garantimos que tem exatamente 11 dígitos
	if len(cpf) != 11 {
		return false
	}

	// CPF inválido conhecido (todos dígitos iguais)
	switch cpf {
	case "00000000000", "11111111111", "22222222222", "33333333333",
		"44444444444", "55555555555", "66666666666", "77777777777",
		"88888888888", "99999999999":
		return false
	}

	var sum int
	for i := range 9 {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	firstDigit := (sum * 10) % 11
	if firstDigit == 10 {
		firstDigit = 0
	}
	if firstDigit != int(cpf[9]-'0') {
		return false
	}

	sum = 0
	for i := range 10 {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	secondDigit := (sum * 10) % 11
	if secondDigit == 10 {
		secondDigit = 0
	}
	if secondDigit != int(cpf[10]-'0') {
		return false
	}

	return true
}

func ValidPhone(phone string) bool {
	return PhoneRX.MatchString(phone)
}

func ValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}

// Função conveniente para validar qualquer struct
func ValidateStruct(ctx context.Context, data any) Evaluator {
	validator := NewStructValidator(data)
	return validator.Valid(ctx)
}

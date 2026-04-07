package transactions

import "strings"

const (
	CategorySalary        = "sueldo"
	CategorySales         = "ventas"
	CategoryTransfer      = "transferencia"
	CategoryFood          = "comida"
	CategoryGroceries     = "supermercado"
	CategoryTransport     = "transporte"
	CategoryFuel          = "combustible"
	CategoryServices      = "servicios"
	CategoryRent          = "alquiler"
	CategoryTaxes         = "impuestos"
	CategoryPharmacy      = "farmacia"
	CategorySubscriptions = "suscripciones"
	CategoryLeisure       = "ocio"
	CategoryClothing      = "ropa"
	CategorySavings       = "ahorro"
	CategoryInvestment    = "inversion"
	CategoryDebts         = "deudas"
	CategoryOther         = "otros"
)

type transactionCategoryDefinition struct {
	Option  TransactionCategoryOption
	Aliases []string
}

var categoryNormalizer = strings.NewReplacer(
	"á", "a",
	"é", "e",
	"í", "i",
	"ó", "o",
	"ú", "u",
	"ü", "u",
	"_", " ",
	"-", " ",
	"/", " ",
)

var transactionCategoryDefinitions = []transactionCategoryDefinition{
	{
		Option: TransactionCategoryOption{
			Key:          CategorySalary,
			Label:        "Sueldo",
			AllowedTypes: []TransactionType{Income},
		},
		Aliases: []string{"salario", "nomina", "nomina mensual"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategorySales,
			Label:        "Ventas",
			AllowedTypes: []TransactionType{Income},
		},
		Aliases: []string{"venta", "cobro", "comisiones", "comision"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryTransfer,
			Label:        "Transferencia",
			AllowedTypes: []TransactionType{Income, Expense},
		},
		Aliases: []string{"transfer", "transferencias"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryFood,
			Label:        "Comida",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"restaurante", "restaurant", "delivery", "cafe", "cafeteria"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryGroceries,
			Label:        "Supermercado",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"super", "almacen", "grocery"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryTransport,
			Label:        "Transporte",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"taxi", "uber", "cabify", "sube", "colectivo", "tren", "peaje", "estacionamiento"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryFuel,
			Label:        "Combustible",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"nafta", "gasoil", "gasolina"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryServices,
			Label:        "Servicios",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"servicio", "internet", "telefono", "luz", "agua", "gas", "seguro", "celular"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryRent,
			Label:        "Alquiler",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"renta", "expensas"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryTaxes,
			Label:        "Impuestos",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"impuesto", "arba", "afip", "monotributo"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryPharmacy,
			Label:        "Farmacia",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"salud", "medico", "medica", "remedios"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategorySubscriptions,
			Label:        "Suscripciones",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"suscripcion", "netflix", "spotify", "prime", "disney"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryLeisure,
			Label:        "Ocio",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"entretenimiento", "salida", "cine", "juegos"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryClothing,
			Label:        "Ropa",
			AllowedTypes: []TransactionType{Expense},
		},
		Aliases: []string{"indumentaria", "vestimenta"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategorySavings,
			Label:        "Ahorro",
			AllowedTypes: []TransactionType{Income, Expense},
		},
		Aliases: []string{"ahorros"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryInvestment,
			Label:        "Inversión",
			AllowedTypes: []TransactionType{Income, Expense},
		},
		Aliases: []string{"inversiones", "invertir", "inversor"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryDebts,
			Label:        "Deudas",
			AllowedTypes: []TransactionType{Income, Expense},
		},
		Aliases: []string{"deuda", "prestamo", "prestamos", "cuota", "cuotas"},
	},
	{
		Option: TransactionCategoryOption{
			Key:          CategoryOther,
			Label:        "Otros",
			AllowedTypes: []TransactionType{Income, Expense},
		},
		Aliases: []string{"otro", "general", "varios", "misc"},
	},
}

func listTransactionCategories() []TransactionCategoryOption {
	options := make([]TransactionCategoryOption, 0, len(transactionCategoryDefinitions))

	for _, definition := range transactionCategoryDefinitions {
		allowedTypes := append([]TransactionType(nil), definition.Option.AllowedTypes...)
		options = append(options, TransactionCategoryOption{
			Key:          definition.Option.Key,
			Label:        definition.Option.Label,
			AllowedTypes: allowedTypes,
		})
	}

	return options
}

func normalizeTransactionCategoryInput(value string, transactionType TransactionType) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", ErrTransactionCategoryRequired
	}

	key, found := resolveTransactionCategoryKey(normalized)
	if !found {
		return "", ErrTransactionCategoryInvalid
	}

	if !categoryAllowsType(key, transactionType) {
		return "", ErrTransactionCategoryTypeMismatch
	}

	return key, nil
}

func normalizeTransactionCategoryFilter(value string) (string, error) {
	key, found := resolveTransactionCategoryKey(value)
	if !found {
		return "", ErrTransactionCategoryInvalid
	}

	return key, nil
}

func normalizeStoredTransactionCategory(transaction Transaction) string {
	if transaction.TransferID != nil {
		return CategoryTransfer
	}

	if isInitialBalanceTransaction(transaction) {
		return CategoryOther
	}

	key, found := resolveTransactionCategoryKey(transaction.Category)
	if !found {
		return CategoryOther
	}

	return key
}

func resolveTransactionCategoryLabel(key string) string {
	for _, definition := range transactionCategoryDefinitions {
		if definition.Option.Key == key {
			return definition.Option.Label
		}
	}

	return "Otros"
}

func resolveTransactionCategoryKey(value string) (string, bool) {
	normalized := normalizeCategoryToken(value)
	if normalized == "" {
		return "", false
	}

	for _, definition := range transactionCategoryDefinitions {
		candidates := buildCategoryCandidates(definition)
		for _, candidate := range candidates {
			if normalized == candidate {
				return definition.Option.Key, true
			}
		}
	}

	for _, definition := range transactionCategoryDefinitions {
		candidates := buildCategoryCandidates(definition)
		for _, candidate := range candidates {
			if candidate != "" && strings.Contains(normalized, candidate) {
				return definition.Option.Key, true
			}
		}
	}

	return "", false
}

func categoryAllowsType(key string, transactionType TransactionType) bool {
	for _, definition := range transactionCategoryDefinitions {
		if definition.Option.Key != key {
			continue
		}

		for _, allowedType := range definition.Option.AllowedTypes {
			if allowedType == transactionType {
				return true
			}
		}

		return false
	}

	return false
}

func buildCategoryCandidates(definition transactionCategoryDefinition) []string {
	candidates := make([]string, 0, len(definition.Aliases)+2)
	candidates = append(candidates, normalizeCategoryToken(definition.Option.Key))
	candidates = append(candidates, normalizeCategoryToken(definition.Option.Label))

	for _, alias := range definition.Aliases {
		candidates = append(candidates, normalizeCategoryToken(alias))
	}

	return candidates
}

func normalizeCategoryToken(value string) string {
	normalized := categoryNormalizer.Replace(strings.ToLower(strings.TrimSpace(value)))
	return strings.Join(strings.Fields(normalized), " ")
}

func isInitialBalanceTransaction(transaction Transaction) bool {
	return normalizeCategoryToken(transaction.Title) == normalizeCategoryToken(InitialBalanceTitle) &&
		normalizeCategoryToken(transaction.Description) == normalizeCategoryToken(InitialBalanceDescription) &&
		transaction.Type == Income
}

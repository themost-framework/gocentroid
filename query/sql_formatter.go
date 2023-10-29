package query

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func FormatMethodName(name string) string {
	re := regexp.MustCompile(`^\$(\w+)$`)
	result := string(re.ReplaceAll([]byte(name), []byte("${1}")))
	result = strings.ToUpper(string(result[:1])) + result[1:]
	return result
}

type SqlDialectOptions struct {
	Space    string
	Select   string
	Update   string
	Insert   string
	Delete   string
	Or       string
	And      string
	OrderBy  string
	GroupBy  string
	Distinct string
	Join     string
	Set      string
	As       string
	Inner    string
	Right    string
	Left     string
	Add      string
	Subtract string
	Multiply string
	Divide   string
	Modulo   string
}

func DefaultSqlDialectOptions() SqlDialectOptions {
	return SqlDialectOptions{
		Space:    " ",
		Select:   "SELECT",
		Update:   "UPDATE",
		Insert:   "INSERT",
		Delete:   "DELETE",
		Or:       "OR",
		And:      "AND",
		OrderBy:  "ORDER BY",
		GroupBy:  "GROUP BY",
		Distinct: "DISTINCT",
		Join:     "JOIN",
		Set:      "SET",
		As:       "AS",
		Inner:    "INNER JOIN",
		Right:    "RIGHT JOIN",
		Left:     "LEFT JOIN",
		Add:      "+",
		Subtract: "-",
		Multiply: "*",
		Divide:   "/",
		Modulo:   "%",
	}
}

type SqlDialect struct {
	NameFormat string
	Options    SqlDialectOptions
}

func (dialect *SqlDialect) Init(options SqlDialectOptions) {
	dialect.Options = options
}

func (dialect *SqlDialect) EscapeName(value string) (string, error) {
	var name = TrimNameReference(value)
	var member = strings.Split(name, ".")
	if len(member) > 2 {
		name = strings.Join(member[len(member)-2:], ".")
	}
	return DefaultNameValidator().Escape(name, dialect.NameFormat)
}

func (dialect *SqlDialect) Escape(value any) (string, error) {
	if value == nil {
		return "NULL", nil
	}
	// format query element
	var thisType = reflect.ValueOf(dialect)
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Map {
		// get keys
		keys := val.MapKeys()
		if len(keys) > 0 {
			// get first key
			key := keys[0].String()
			// and check if it starts with "$" which indicates that it's a query expression
			if strings.HasPrefix(key, "$") {
				// use dialect to convert value
				var method = thisType.MethodByName(FormatMethodName(key))
				if !method.IsValid() {
					return "", errors.New("the specified method is unknown")
				}
				// try to get arguments
				args := val.MapIndex(keys[0])
				valueOrValues := args.Elem()
				arguments := []reflect.Value{}
				if valueOrValues.Kind() == reflect.Slice {
					// append arguments
					for i := 0; i < valueOrValues.Len(); i++ {
						var argument = valueOrValues.Index(i)
						arguments = append(arguments, argument)
					}
				} else {
					// append single argument
					arguments = append(arguments, valueOrValues)
				}
				// call dialect method
				res := method.Call(arguments)
				// get result and error
				var result = res[0]
				var err = res[1]
				if !err.IsNil() {
					return "", errors.New(err.String())
				}
				return result.String(), nil

			}
		}
	}
	return (&SqlUtils{}).Escape(value), nil
}

func (dialect *SqlDialect) EscapeLeftRightOperand(left any, right any) (string, string, error) {
	// escape left operation
	leftOperand, err := dialect.Escape(left)
	if err != nil {
		return "", "", err
	}
	rightOperand, err := dialect.Escape(right)
	if err != nil {
		return "", "", err
	}
	return leftOperand, rightOperand, nil
}

func (dialect *SqlDialect) Name(name string) (string, error) {
	return dialect.EscapeName(name)
}

func (dialect *SqlDialect) GetField(field string) (string, error) {
	return dialect.EscapeName(field)
}

func (dialect *SqlDialect) Eq(left any, right any) (string, error) {
	leftOperand, rightOperand, err := dialect.EscapeLeftRightOperand(left, right)
	if err != nil {
		return "", err
	}
	if leftOperand == "NULL" {
		return "", errors.New("left operand cannot be null")
	}
	if rightOperand == "NULL" {
		return fmt.Sprintf("%s IS NULL", leftOperand), nil
	}
	return fmt.Sprintf("%s = %s", leftOperand, rightOperand), nil
}

func (dialect *SqlDialect) Ne(left any, right any) (string, error) {
	leftOperand, rightOperand, err := dialect.EscapeLeftRightOperand(left, right)
	if err != nil {
		return "", err
	}
	if leftOperand == "NULL" {
		return "", errors.New("left operand cannot be null")
	}
	if rightOperand == "NULL" {
		return fmt.Sprintf("NOT %s IS NULL", leftOperand), nil
	}
	return fmt.Sprintf("%s <> %s", leftOperand, rightOperand), nil
}

func (dialect *SqlDialect) Gt(left any, right any) (string, error) {
	leftOperand, rightOperand, err := dialect.EscapeLeftRightOperand(left, right)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s > %s", leftOperand, rightOperand), nil
}

func (dialect *SqlDialect) Gte(left any, right any) (string, error) {
	leftOperand, rightOperand, err := dialect.EscapeLeftRightOperand(left, right)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s >= %s", leftOperand, rightOperand), nil
}

func (dialect *SqlDialect) Lt(left any, right any) (string, error) {
	leftOperand, rightOperand, err := dialect.EscapeLeftRightOperand(left, right)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s < %s", leftOperand, rightOperand), nil
}

func (dialect *SqlDialect) Lte(left any, right any) (string, error) {
	leftOperand, rightOperand, err := dialect.EscapeLeftRightOperand(left, right)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s <= %s", leftOperand, rightOperand), nil
}

func (dialect *SqlDialect) EscapeBinaryExpr(operator string) func(args ...any) (string, error) {
	return func(args ...any) (string, error) {
		var exprs = []string{}
		for _, arg := range args {
			expr, err := dialect.Escape(arg)
			if err != nil {
				return "", err
			}
			exprs = append(exprs, expr)
		}
		return strings.Join(exprs, dialect.Options.Space+operator+dialect.Options.Space), nil
	}
}

func (dialect *SqlDialect) Or(args ...any) (string, error) {
	return dialect.EscapeBinaryExpr(dialect.Options.Or)(args...)
}

func (dialect *SqlDialect) And(args ...any) (string, error) {
	return dialect.EscapeBinaryExpr(dialect.Options.And)(args...)
}

func (dialect *SqlDialect) Add(args ...any) (string, error) {
	return dialect.EscapeBinaryExpr(dialect.Options.Add)(args...)
}

func (dialect *SqlDialect) Subtract(args ...any) (string, error) {
	return dialect.EscapeBinaryExpr(dialect.Options.Subtract)(args...)
}

func (dialect *SqlDialect) Multiply(args ...any) (string, error) {
	return dialect.EscapeBinaryExpr(dialect.Options.Multiply)(args...)
}

func (dialect *SqlDialect) Mul(args ...any) (string, error) {
	return dialect.Multiply(args...)
}

func (dialect *SqlDialect) Divide(args ...any) (string, error) {
	return dialect.EscapeBinaryExpr(dialect.Options.Divide)(args...)
}

func (dialect *SqlDialect) Div(args ...any) (string, error) {
	return dialect.Divide(args...)
}

func (dialect *SqlDialect) Count(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("COUNT(%s)", expr), nil
}

func (dialect *SqlDialect) Min(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("MIN(%s)", expr), nil
}

func (dialect *SqlDialect) Max(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("MIN(%s)", expr), nil
}

func (dialect *SqlDialect) Sum(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SUM(%s)", expr), nil
}

func (dialect *SqlDialect) Average(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("AVG(%s)", expr), nil
}

func (dialect *SqlDialect) Avg(arg any) (string, error) {
	return dialect.Average(arg)
}

func (dialect *SqlDialect) Length(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("LENGTH(%s)", expr), nil
}

func (dialect *SqlDialect) Trim(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("TRIM(%s)", expr), nil
}

func (dialect *SqlDialect) ToUpper(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("UPPER(%s)", expr), nil
}

func (dialect *SqlDialect) ToLower(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("LOWER(%s)", expr), nil
}

func (dialect *SqlDialect) Year(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("YEAR(%s)", expr), nil
}

func (dialect *SqlDialect) Month(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("MONTH(%s)", expr), nil
}

func (dialect *SqlDialect) DayOfMonth(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("DAY(%s)", expr), nil
}

func (dialect *SqlDialect) Hour(arg ...any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("HOUR(%s)", expr), nil
}

func (dialect *SqlDialect) Minute(arg ...any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("MINUTE(%s)", expr), nil
}

func (dialect *SqlDialect) Second(arg any) (string, error) {
	expr, err := dialect.Escape(arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SECOND(%s)", expr), nil
}

func (dialect *SqlDialect) Modulo(args ...any) (string, error) {
	var exprs = []string{}
	for _, arg := range args {
		expr, err := dialect.Escape(arg)
		if err != nil {
			return "", err
		}
		exprs = append(exprs, expr)
	}
	result := "("
	result += strings.Join(exprs, dialect.Options.Space+dialect.Options.Modulo+dialect.Options.Space)
	result += ")"
	return result, nil
}

func (dialect *SqlDialect) Mod(args ...any) (string, error) {
	return dialect.Modulo(args...)
}

func (dialect *SqlDialect) Cond(ifExpr any, thenExpr any, elseExpr any) (string, error) {
	arg0, err := dialect.Escape(ifExpr)
	if err != nil {
		return "", err
	}
	arg1, err := dialect.Escape(thenExpr)
	if err != nil {
		return "", err
	}
	arg2, err := dialect.Escape(thenExpr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(CASE %s WHEN 1 THEN %s ELSE %s END)", arg0, arg1, arg2), nil
}

type SqlFormatter struct {
}

func (formatter *SqlFormatter) FormatSelect(query QueryExpression) *string {
	return nil
}

package unit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	alphaRegexString               = "^[a-zA-Z]+$"
	alphaNumericRegexString        = "^[a-zA-Z0-9]+$"
	alphaUnicodeRegexString        = "^[\\p{L}]+$"
	alphaUnicodeNumericRegexString = "^[\\p{L}\\p{N}]+$"
	numericRegexString             = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	numberRegexString              = "^[0-9]+$"
	hexadecimalRegexString         = "^[0-9a-fA-F]+$"
	hexcolorRegexString            = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
	rgbRegexString                 = "^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"
	rgbaRegexString                = "^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	hslRegexString                 = "^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*\\)$"
	hslaRegexString                = "^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	emailRegexString               = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	base64RegexString              = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	iSBN10RegexString              = "^(?:[0-9]{9}X|[0-9]{10})$"
	iSBN13RegexString              = "^(?:(?:97(?:8|9))[0-9]{10})$"
	uUID3RegexString               = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uUID4RegexString               = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUID5RegexString               = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUIDRegexString                = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	aSCIIRegexString               = "^[\x00-\x7F]*$"
	printableASCIIRegexString      = "^[\x20-\x7E]*$"
	multibyteRegexString           = "[^\x00-\x7F]"
	dataURIRegexString             = "^data:.+\\/(.+);base64$"
	latitudeRegexString            = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	longitudeRegexString           = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	sSNRegexString                 = `^\d{3}[- ]?\d{2}[- ]?\d{4}$`
)

var (
	alphaRegex               = regexp.MustCompile(alphaRegexString)
	alphaNumericRegex        = regexp.MustCompile(alphaNumericRegexString)
	alphaUnicodeRegex        = regexp.MustCompile(alphaUnicodeRegexString)
	alphaUnicodeNumericRegex = regexp.MustCompile(alphaUnicodeNumericRegexString)
	numericRegex             = regexp.MustCompile(numericRegexString)
	numberRegex              = regexp.MustCompile(numberRegexString)
	hexadecimalRegex         = regexp.MustCompile(hexadecimalRegexString)
	hexcolorRegex            = regexp.MustCompile(hexcolorRegexString)
	rgbRegex                 = regexp.MustCompile(rgbRegexString)
	rgbaRegex                = regexp.MustCompile(rgbaRegexString)
	hslRegex                 = regexp.MustCompile(hslRegexString)
	hslaRegex                = regexp.MustCompile(hslaRegexString)
	emailRegex               = regexp.MustCompile(emailRegexString)
	base64Regex              = regexp.MustCompile(base64RegexString)
	iSBN10Regex              = regexp.MustCompile(iSBN10RegexString)
	iSBN13Regex              = regexp.MustCompile(iSBN13RegexString)
	uUID3Regex               = regexp.MustCompile(uUID3RegexString)
	uUID4Regex               = regexp.MustCompile(uUID4RegexString)
	uUID5Regex               = regexp.MustCompile(uUID5RegexString)
	uUIDRegex                = regexp.MustCompile(uUIDRegexString)
	aSCIIRegex               = regexp.MustCompile(aSCIIRegexString)
	printableASCIIRegex      = regexp.MustCompile(printableASCIIRegexString)
	multibyteRegex           = regexp.MustCompile(multibyteRegexString)
	dataURIRegex             = regexp.MustCompile(dataURIRegexString)
	latitudeRegex            = regexp.MustCompile(latitudeRegexString)
	longitudeRegex           = regexp.MustCompile(longitudeRegexString)
	sSNRegex                 = regexp.MustCompile(sSNRegexString)
)

// M map type string of interfaces
type M map[string]interface{}

type checkEmailResponse struct {
	Email       string  `json:"email"`
	DidYouMean  string  `json:"did_you_mean"`
	User        string  `json:"user"`
	Domain      string  `json:"domain"`
	FormatValid bool    `json:"format_valid"`
	MXFound     bool    `json:"mx_found"`
	SMTPCheck   bool    `json:"smtp_check"`
	Role        bool    `json:"role"`
	Disposable  bool    `json:"disposable"`
	Free        bool    `json:"free"`
	Score       float32 `json:"score"`
}

func isRealEmail(email string) bool {
	if emailRegex.MatchString(email) {
		return true
	}

	return false
}

// EX
//
//	errs = validate(map[string]interface{}{
//		"username":       username,
//		"password":       password,
//		"password_check": passwordCheck,
//		"email":          email,
//		"secret_pin":     secretPIN,
//	}, map[string]string{
//
//		"username":       "required|min:6|max:15|alphanum|unique:users",
//		"password":       "required|min:6|max:15|alphanum",
//		"password_check": "required|same:password",
//		"email":          "required|email|unique:users",
//		"secret_pin":     "required|len:6",
//	})
func Validate(values map[string]interface{}, rules map[string]string) M {
	var errorBag = make(M)

	for key, value := range values {
		var valRules = strings.Split(rules[key], "|")

		for _, rule := range valRules {
			var ruleData = strings.Split(rule, ":")

			// required rule
			if ruleData[0] == "required" {
				if value == "" {
					errorBag[key] = fmt.Sprintf("The %s field is required.", key)
					break
				}
			}

			// email rule
			if ruleData[0] == "email" {
				if emailRegex.MatchString(value.(string)) == false {
					errorBag[key] = fmt.Sprintf("The %s must be a valid email address.", key)
					break
				}
			}

			// unique rule opts: tableName
			if ruleData[0] == "unique" {
				// var tableName = ruleData[1]
				// if models.CheckIfExist(tableName, key, value) {
				// 	errorBag[key] = fmt.Sprintf("The %s has already been taken.", key)
				// 	break
				// }
			}

			// exists rule opts: tableName
			if ruleData[0] == "exists" {
				// var tableName = ruleData[1]
				// if !models.CheckIfExist(tableName, key, value) {
				// 	errorBag[key] = fmt.Sprintf("The %s does not exist.", key)
				// 	break
				// }
			}

			// alphanum rule
			if ruleData[0] == "alphanum" {
				if alphaNumericRegex.MatchString(value.(string)) == false {
					errorBag[key] = fmt.Sprintf("The %s may only contain letters and numbers.", key)
					break
				}
			}

			// integer rule
			if ruleData[0] == "integer" {
				switch v := value.(type) {
				case int:
					break
				case string:
					if numberRegex.MatchString(v) == false {
						errorBag[key] = fmt.Sprintf("The %s must be an integer.", key)
						break
					}
				}
			}

			// min rule
			if ruleData[0] == "min" {
				ruleVal, err := strconv.Atoi(ruleData[1])
				if err == nil {
					switch v := value.(type) {
					case int:
						if v < ruleVal {
							errorBag[key] = fmt.Sprintf("The %s must be at least %d.", key, ruleVal)
							break
						}
					case float64:
						var floatRuleVal float64
						floatRuleVal, err = strconv.ParseFloat(ruleData[1], 64)
						if err == nil {
							if v < floatRuleVal {
								errorBag[key] = fmt.Sprintf("The %s must be at least %.2f.", key, floatRuleVal)
								break
							}
						}
					case string:
						if len(v) < ruleVal {
							errorBag[key] = fmt.Sprintf("The %s must be at least %d characters.", key, ruleVal)
							break
						}
					}
				}
			}

			// max rule
			if ruleData[0] == "max" {
				switch v := value.(type) {
				case int:
					ruleVal, err := strconv.Atoi(ruleData[1])
					if err == nil {
						if v > ruleVal {
							errorBag[key] = fmt.Sprintf("The %s may not be greater than %d.", key, ruleVal)
							break
						}
					}
				case float64:
					floatRuleVal, err := strconv.ParseFloat(ruleData[1], 64)
					if err == nil {
						if v > floatRuleVal {
							errorBag[key] = fmt.Sprintf("The %s may not be greater than %.2f.", key, floatRuleVal)
							break
						}
					}
				case string:
					ruleVal, err := strconv.Atoi(ruleData[1])
					if err == nil {
						if len(v) > ruleVal {
							errorBag[key] = fmt.Sprintf("The %s may not be greater than %d characters.", key, ruleVal)
							break
						}
					}
				}
			}

			// len rule
			if ruleData[0] == "len" {
				ruleVal, err := strconv.Atoi(ruleData[1])
				if err == nil {
					if len(value.(string)) != ruleVal {
						errorBag[key] = fmt.Sprintf("The %s must have exactly %d characters.", key, ruleVal)
						break
					}
				}
			}

			// same rule
			if ruleData[0] == "same" {
				var otherField = ruleData[1]
				if value != values[otherField] {
					errorBag[key] = fmt.Sprintf("The %s and %s must match.", key, otherField)
					break
				}
			}

			// decimal rule
			if ruleData[0] == "decimal" {
				switch v := value.(type) {
				case float32, float64:
					break
				case string:
					if numericRegex.MatchString(v) == false {
						errorBag[key] = fmt.Sprintf("The %s must be a decimal number.", key)
						break
					}
				}
			}

			// accepted rule
			if ruleData[0] == "accepted" {
				switch v := value.(type) {
				case bool:
					if !v {
						errorBag[key] = fmt.Sprintf("The %s must be accepted by checking the box.", key)
						break
					}
				case string:
					if v == "" {
						errorBag[key] = fmt.Sprintf("The %s must be accepted by checking the box.", key)
						break
					}
				}
			}

			// phone rule
			if ruleData[0] == "phone" {
				// switch v := value.(type) {
				// case string:
				// 	num, err := libphonenumber.Parse(v, "")
				// 	if err != nil {
				// 		errorBag[key] = fmt.Sprintf("Mobile number: %s.", err.Error())
				// 		break
				// 	}
				// 	if !libphonenumber.IsValidNumber(num) {
				// 		errorBag[key] = fmt.Sprintf("The %s must be a valid mobile phone number.", key)
				// 		break
				// 	}
				// }
			}

			// in rule
			if ruleData[0] == "in" {
				var selection = strings.Split(ruleData[1], ",")
				var ruleMatched bool

				switch v := value.(type) {
				case string:
					if v != "" {
						for _, str := range selection {
							if v == str {
								ruleMatched = true
								break
							}
						}

						if !ruleMatched {
							errorBag[key] = fmt.Sprintf("The selected %s is invalid.", key)
						}
					}
				}
			}

		}
	}

	if len(errorBag) > 0 {
		return errorBag
	}

	return nil
}

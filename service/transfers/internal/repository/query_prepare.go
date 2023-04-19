package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model"
	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
)

func prepQuery(f treegrid_model.FilterParams) (map[string]string, map[string][]interface{}) {
	FilterWhere := map[string]string{}
	FilterArgs := map[string][]interface{}{}

	// Filter process
	var curField, curFieldValue, curOperation, curMarker string

	for key, el := range f.MainFilter() {
		if key == "id" || strings.Contains(key, "Filter") {
			continue
		}

		// Check if cur field is child's or parent's and generate preWhere and postWhere correspondingly
		if model.ItemsFields[key] != "" {
			curMarker = "child"
		} else {
			// cur column is a parent column
			curMarker = "parent"
		}

		curField = key
		curOperation = f.MainFilter()[curField+"Filter"].(string)
		curFieldValue = el.(string)

		// for child item
		if model.ItemsFields[key] != "" {
			curField = model.ItemsFields[key]
		}

		if model.FieldAliases[key] != "" {
			curField = model.FieldAliases[key]
			if curField[:11] == "STR_TO_DATE" {
				curFieldValue = "STR_TO_DATE('" + el.(string) + "','%m/%d/%Y')"
				FilterWhere[curMarker] += " AND " + curField + model.FieldAliasesDate[curOperation] + curFieldValue
				continue
			}
		}

		if curOperation != "" {
			switch curOperation {
			case "1":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " = ? ", "")
			case "2":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " != ? ", "")
			case "3":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " < ? ", "")
			case "4":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " <= ? ", "")
			case "5":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " > ? ", "")
			case "6":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " >= ? ", "")
			case "7":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " LIKE ? ", "end")
			case "8":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " NOT LIKE ? ", "end")
			case "9":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " LIKE ? ", "start")
			case "10":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " NOT LIKE ? ", "start")
			case "11":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " LIKE ? ", "both")
			case "12":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " NOT LIKE ? ", "both")
			case "13":
				prepareFilterbyDeli(curMarker, curField, curFieldValue, FilterWhere, FilterArgs, " AND ", " IN ? ", "")
			}
		}
	}

	return FilterWhere, FilterArgs
}

func prepareFilterbyDeli(
	curMarker string, curField string, curFieldValue string,
	filterWhere map[string]string, filterArgs map[string][]interface{},
	condition string, conditionVal string, modPosition string) {
	// Check if filter have any delimited values
	splitValue := strings.Split(curFieldValue, ";")
	// Check if filter have any Range values
	rangeValue := strings.Split(curFieldValue, "~")
	inValues := strings.Split(curFieldValue, ",")

	// condition = " AND ", conditionVal = " LIKE ? ", modPosition = "end"
	if len(splitValue) > 1 {
		for i := 0; i < len(splitValue); i++ {
			if i == 0 {
				filterWhere[curMarker] += condition + curField + conditionVal
			} else {
				filterWhere[curMarker] += " OR " + curField + conditionVal
			}

			if modPosition == "both" {
				filterArgs[curMarker] = append(filterArgs[curMarker], "%"+splitValue[i]+"%")
			} else if modPosition == "start" {
				filterArgs[curMarker] = append(filterArgs[curMarker], "%"+splitValue[i])
			} else if modPosition == "end" {
				filterArgs[curMarker] = append(filterArgs[curMarker], splitValue[i]+"%")
			} else {
				filterArgs[curMarker] = append(filterArgs[curMarker], splitValue[i])
			}
		}
		return
	}

	if len(rangeValue) > 1 {
		if !utils.IsDateValue(rangeValue[0]) {
			// Check if Range is numeric value
			start, _ := strconv.Atoi(rangeValue[0])
			end, _ := strconv.Atoi(rangeValue[1])
			count := 0
			for i := start; i < end; i++ {
				if count == 0 {
					filterWhere[curMarker] += condition + curField + conditionVal
				} else {
					filterWhere[curMarker] += " OR " + curField + conditionVal
				}
				if modPosition == "both" {
					filterArgs[curMarker] = append(filterArgs[curMarker], "%"+strconv.Itoa(i)+"%")
				} else if modPosition == "start" {
					filterArgs[curMarker] = append(filterArgs[curMarker], "%"+strconv.Itoa(i))
				} else if modPosition == "end" {
					filterArgs[curMarker] = append(filterArgs[curMarker], strconv.Itoa(i)+"%")
				} else {
					filterArgs[curMarker] = append(filterArgs[curMarker], strconv.Itoa(i))
				}
				count += 1
			}
		} else if utils.IsDateValue(rangeValue[0]) {
			// Check if Range is date value
			var err error
			start, err := time.Parse("01/02/2006", splitValue[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "err3", err)
				log.Fatal(err)
			}
			end, err := time.Parse("01/02/2006", splitValue[1])
			count := 0
			if err == nil {
				for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
					if count == 0 {
						filterWhere[curMarker] += condition + curField + conditionVal
					} else {
						filterWhere[curMarker] += " OR " + curField + conditionVal
					}
					if modPosition == "both" {
						filterArgs[curMarker] = append(filterArgs[curMarker], "%"+d.Format("2006-01-02")+"%")
					} else if modPosition == "start" {
						filterArgs[curMarker] = append(filterArgs[curMarker], "%"+d.Format("2006-01-02"))
					} else if modPosition == "end" {
						filterArgs[curMarker] = append(filterArgs[curMarker], d.Format("2006-01-02")+"%")
					} else {
						filterArgs[curMarker] = append(filterArgs[curMarker], d.Format("2006-01-02"))
					}
					count += 1
				}
			} else {
				fmt.Fprintln(os.Stderr, "err4", err)
				log.Fatal(err)
			}
		} else {
			filterWhere[curMarker] += condition + curField + conditionVal
			// filterArgs[curMarker] = append(filterArgs[curMarker], "%"+curFieldValue+"%")
			filterArgs[curMarker] = append(filterArgs[curMarker], curFieldValue)
		}
		return
	}
	// condition , conditionVal, modPosition
	// " AND ", " IN ? ", ""
	if len(inValues) > 1 {
		noOfValues := ""
		for i := 0; i < len(inValues); i++ {
			noOfValues += "?,"
			filterArgs[curMarker] = append(filterArgs[curMarker], inValues[i])
		}
		noOfValues = noOfValues[:len(noOfValues)-1]
		filterWhere[curMarker] += condition + curField + " IN (" + noOfValues + ")"

		return
	}

	filterWhere[curMarker] += condition + curField + conditionVal
	// filterArgs[curMarker] = append(filterArgs[curMarker], "%"+curFieldValue+"%")
	filterArgs[curMarker] = append(filterArgs[curMarker], curFieldValue)
}

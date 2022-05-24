package extract

import (
	"fmt"
	"bufio"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	
	"github.com/InstIDEA/ddjj/parser/declaration"
)

var totalAssets int64

var assetsItemNumber int

var skipAssets = []string{
	"#",
	"DESCRIPCIÓN",
	"EMPRESA",
	"RUC",
	"PAÍS",
	"CANT.",
	"PRECIO UNI.",
	"IMPORTE",
}

// Assets returns other assets owned by the official.
func Assets(scanner *bufio.Scanner) ([]*declaration.OtherAsset, error) {
	scanner = MoveUntil2NL(scanner, "1.9 OTROS ACTIVOS", false)

	// Also wants to skip item number
	assetsItemNumber = 1
	skipAssets = append(skipAssets, strconv.Itoa(assetsItemNumber))

	var assets []*declaration.OtherAsset

	values, nextPage := getAssetValues(scanner, 0, false, 6)
	for values[0] != "" {
		asset := getAsset(scanner, values)
		assets = append(assets, asset...)
	
		if nextPage {
			assetsItemNumber = 1
		} else {
			assetsItemNumber++
		}
		// Also wants to skip item number
		skipAssets[len(skipAssets)-1] = strconv.Itoa(assetsItemNumber)

		values, nextPage = getAssetValues(scanner, 0, false, 6)
	}

	total := addAssets(assets)
	if total == 0 {
		return nil, errors.New("failed when extracting other assets")
	}

	if total != totalAssets {
		return nil, errors.New("other assets do not match")
	}

	// Reset variables for next call.
	totalAssets = 0
	assetsItemNumber = 0

	return assets, nil
}

func getAssetValues(scanner *bufio.Scanner, index int, remaining bool, lim int) (values [7]string, nextPage bool) {
	line, _ := getAssetLine(scanner)
	for line != "" {
		values[index] = line
		// After reading all the possible values for a single item.
		if index == lim {
			return
		}

		index++

		line, nextPage = getAssetLine(scanner)
	}

	if remaining {
		return
	}

	return [7]string{}, false
}

func getAsset(scanner *bufio.Scanner, values [7]string) []*declaration.OtherAsset {
	// En algunos casos, el importe del primer activo está al final de la lista
	// de activos. Por ejemplo Juan Afara 2014
	if !isNumber(values[6]) {
		return getAsset2(scanner, values)
	}

	return []*declaration.OtherAsset{getAsset1(values)}
}

func getAsset1(values [7]string) *declaration.OtherAsset {
	if !isNumber(values[6]){
		return &declaration.OtherAsset{
			Descripcion: values[0],
			Empresa:     values[1],
			RUC:         values[2],
			Pais:        values[3],
			Cantidad:    stringToInt64("1"),
			Precio:      stringToInt64(values[4]),
			Importe:     stringToInt64(values[5]),
		}
	}
	return &declaration.OtherAsset{
		Descripcion: values[0],
		Empresa:     values[1],
		RUC:         values[2],
		Pais:        values[3],
		Cantidad:    stringToInt64(values[4]),
		Precio:      stringToInt64(values[5]),
		Importe:     stringToInt64(values[6]),
	}
}

func getAsset2(scanner *bufio.Scanner, values [7]string) []*declaration.OtherAsset {
	assets := []*declaration.OtherAsset{}

	firstAsset := getAsset1(values)
	assets = append(assets, firstAsset)

	assetsItemNumber++
	skipAssets = append(skipAssets, strconv.Itoa(assetsItemNumber))

	// values[6] is the descripcion in the second element.
	tmp := values[6]
	values, _ = getAssetValues(scanner, 1, false, 5)
	values[0] = tmp
	secondAsset := getAsset1(values)
	assets = append(assets, secondAsset)
	fmt.Println(values)

	// Skip next item number.
	assetsItemNumber++
	skipAssets = append(skipAssets, strconv.Itoa(assetsItemNumber))

	values, nextPage := getAssetValues(scanner, 0, true, 5)
	counter := 0
	for values[1] != "" && !nextPage {
		//assets = append(assets, getAsset(scanner, values)[assetsItemNumber-1])
		assets = append(assets, getAsset1(values))
		//getAsset(scanner,values)

		assetsItemNumber++
		skipAssets = append(skipAssets, strconv.Itoa(assetsItemNumber))
		counter++

		values, nextPage = getAssetValues(scanner, 0, true, 5)
	}

	// The last value is the importe for the first item.
	//firstAsset.Importe = stringToInt64(values[0])

	// Restore skip assets to default state. The caller would remove the other
	// remaining value.
	skipAssets = skipAssets[:len(skipAssets)-counter-2]
	assetsItemNumber = 1

	return assets
}

func getAssetLine(scanner *bufio.Scanner) (line string, nextPage bool) {
	for scanner.Scan() {

		line := strings.ReplaceAll(scanner.Text(), "\n", " ")
		if line == "TOTAL OTROS ACTIVOS" {
			scanner.Scan()
			totalAssets = stringToInt64(strings.ReplaceAll(scanner.Text(), ",", ""))

			// Next page or end.
			scanner = MoveUntil2NL(scanner, "TIPO MUEBLES", true)
			line = scanner.Text()
			if line == ""{
				nextPage = false
			} //else {
			// 	nextPage = true
			// }

			assetsItemNumber = 1
			skipAssets[len(skipAssets)-1] = strconv.Itoa(assetsItemNumber)
		}

		if strings.Contains(line, "OBS:") || strings.Contains(line, "RECEPCIONADO EL:") || strings.Contains(line, "JURADA AL") || strings.Contains(line, "RECEPCIONADO") || strings.Contains(line, "BAJA DEL CARGO"){
			continue
		}
		if isDate(line) || isBarCode(line) {
			continue
		}
		if line == "" || contains(skipAssets, line) {
			continue
		}

		return line, nextPage
	}

	return "", false
}

func addAssets(assets []*declaration.OtherAsset) int64 {
	var total int64
	for _, a := range assets {
		total += a.Importe
	}
	return total
}

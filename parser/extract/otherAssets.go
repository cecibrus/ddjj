package extract

import (
	"fmt"
	"strings"
	"github.com/InstIDEA/ddjj/parser/declaration"
)

func otherAssets(e *Extractor, parser *ParserData) []*declaration.OtherAsset{

	//e.Scanner.Split(Split2NL)

	e.BindFlag(EXTRACTOR_FLAG_1)
	e.BindFlag(EXTRACTOR_FLAG_2)
	//e.BindFlag(EXTRACTOR_FLAG_3)

	var assets []*declaration.OtherAsset
	assets = countAssets(e, assets)
	//var successful int

	//fmt.Println("Tiene esta cantidad de otros activos:", counter)

	e.Rewind()

	fmt.Println(assets)

	return assets

}

// func getJobTitle(e *Extractor) string {

// 	if strings.Contains(e.CurrToken, "CARGO") {
// 		val, check := isKeyValuePair(e.CurrToken, "CARGO")
// 		if check {
// 			return val
// 		}
// 	}

// 	if strings.Contains(e.PrevToken, "CARGO") &&
// 	strings.Contains(e.CurrToken, "FECHA EGRESO") {
// 		return e.NextToken

// 	}

// 	return ""
// }

// func getJobInst(e *Extractor) string {

// 	if strings.Contains(e.PrevToken, "INSTITUCIÓN") &&
// 	strings.Contains(e.NextToken, "ACTO ADM. COM") {
// 		return e.CurrToken
// 	}

// 	if strings.Contains(e.PrevToken, "DIRECCIÓN") &&
// 	isNumber(e.CurrToken) {
// 		return e.NextToken
// 	}

// 	return ""
// }

func countAssets(e *Extractor, assets []*declaration.OtherAsset) []*declaration.OtherAsset {
	asset := &declaration.OtherAsset{ }
	// var counterAcciones int
	// var counterInversiones int
	// var counterCDA int
	// var counterBonos int
	// var counterPatentes int
	// var counterOtros int

	for e.Scan() {
		if strings.Contains(e.CurrToken, "ACCIONES"){
			//counterAcciones++
			//fmt.Println(e.CurrToken)
			values := tokenize(e.CurrToken, 5)
			//asset is added only if it has all of the needed values
			if len(values) == 8 {
				asset = getAsset3(values)
				assets = append(assets, asset)
			} else {
				continue
			}
		} else if strings.Contains(e.CurrToken, "INVERSIONES") {
			//counterInversiones++
			//fmt.Println("Encontro inversiones", counterInversiones)
			values := tokenize(e.CurrToken, 5)
			if len(values) == 8 {
				asset = getAsset3(values)
				assets = append(assets, asset)
			} else {
				continue
			}
		} else if strings.Contains(e.CurrToken, "CERTIFICADO DE DEPOSITOS DE"){
			//counterCDA++
			//fmt.Println("Encontro CDA", counterCDA)
			fmt.Println(e.NextToken)
			// values := tokenize(e.CurrToken, 5)
			// if len(values) == 8 {
			// 	asset = getAsset3(values)
			// 	assets = append(assets, asset)
			// } else {
			// 	continue
			// }
		} else if strings.Contains(e.CurrToken, "BONOS"){
			//counterBonos++
			//fmt.Println(e.CurrToken)
			values := tokenize(e.CurrToken, 5)
			if len(values) == 8 {
				asset = getAsset3(values)
				assets = append(assets, asset)
			} else {
				continue
			}
		} else if strings.Contains(e.CurrToken, "PATENTES"){
			//counterPatentes++
			values := tokenize(e.CurrToken, 5)
			if len(values) == 8 {
				asset = getAsset3(values)
				assets = append(assets, asset)
			} else {
				continue
			}
		}else if strings.Contains(e.CurrToken, "OTROS") && isNumber(e.PrevToken) {
			//counterOtros++
			//fmt.Println("Encontro otros", counterOtros)
			values := tokenize(e.CurrToken, 5)
			fmt.Println(e.CurrToken)
			if len(values) == 8 {
				asset = getAsset3(values)
				assets = append(assets, asset)
			} else {
				continue
			}
		}else{
			continue
		}
	}

	return assets
}

func isAssetFormField(s string) bool {
	formField := []string {
		"#",
		"DESCRIPCION",
		"EMPRESA",
		"RUC",
		"PAIS",
		"CANT.",
		"PRECIO UNI.",
		"IMPORTE",
	}

	s = removeAccents(s)
	for _, value := range formField {
		if isCurrLine(s, value) {
			return true
		}
	}

	return false
}

func getAsset3(values []string) *declaration.OtherAsset {
	return &declaration.OtherAsset{
		Descripcion: values[1],
		Empresa:     values[2],
		RUC:         values[3],
		Pais:        values[4],
		Cantidad:    stringToInt64(values[5]),
		Precio:      stringToInt64(values[6]),
		Importe:     stringToInt64(values[7]),
	}
}
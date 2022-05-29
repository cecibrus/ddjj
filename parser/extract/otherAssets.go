package extract

import (
	"fmt"
	"strings"
	"github.com/pkg/errors"
	"github.com/InstIDEA/ddjj/parser/declaration"
)

func otherAssets(e *Extractor, parser *ParserData) ([]*declaration.OtherAsset, error){

	e.BindFlag(EXTRACTOR_FLAG_1)
	e.BindFlag(EXTRACTOR_FLAG_2)

	var assets []*declaration.OtherAsset
	assets = countAssets(e, assets)

	total := addAssets(assets)
	if total == 0 {
		return nil, errors.New("failed when extracting other assets")
	}

	return assets, nil
}

func countAssets(e *Extractor, assets []*declaration.OtherAsset) []*declaration.OtherAsset {
	asset := &declaration.OtherAsset{ }
	for e.Scan() {
		if strings.Contains(e.CurrToken, "ACCIONES"){
			//Cuando el nombre de la empresa tiene dos lineas, queda en PrevToken, sino queda OBS N/A o el form field y es el caso "normal"
			if strings.Contains(e.PrevToken, "OBS: N/A") || isAssetFormField(e.PrevToken){
				values := tokenize(e.CurrToken, 5)
			//asset is added only if it has all of the needed values
				if len(values) == 8 {
					asset = getAsset3(values)
					assets = append(assets, asset)
				} else {
					continue
				}
			}else{
				//Cuando hay dos lineas, la primera linea queda en PrevToken, en CurrToken queda el indice y ACCIONES y en NextToken la segunda linea del nombre
				//entonces se concatenan PrevToken y NextToken, y luego se vuelve a scanear para tener el resto de los datos
				//Tambien se tiene el caso en el que el nombre se encuentra una linea mas arriba a pesar de no tener dos lineas (Ejemplo: Cartes 2018, Consignataria de Ganado S.A>)
				var name string
				if !strings.Contains(e.NextToken, "OBS: N/A"){
					name = e.PrevToken + " " + e.NextToken
					fmt.Println("Nombre completo: ", name)
				}
				for i := 1; i < 3; i++ {
					e.Scan()
				}
				fmt.Println(e.CurrToken)
				//aditional values that are not in the line but are needed to have the full asset, included the name
				additional := []string{"#","ACCIONES", name}
				values := append(additional, tokenize(e.CurrToken, 4)...)
				if len(values) == 8 {
					asset = getAsset3(values)
					assets = append(assets, asset)
				} else {
					continue
				}

			}
		} else if strings.Contains(e.CurrToken, "CERTIFICADO DE DEPOSITOS DE"){
			//subsequent scans are needed due to the document format
			for i := 1; i < 4; i++ {
				e.Scan()
			}
			//fixed values that are not in the line but are needed to have the full asset
			fixed := []string{"#","CERTIFICADO DE DEPOSITOS DE AHORROS"}
			values := append(fixed, tokenize(e.CurrToken, 4)...)
			asset = getAsset3(values)
			assets = append(assets, asset)
		} else if strings.Contains(e.CurrToken, "INVERSIONES") || strings.Contains(e.CurrToken, "BONOS") || strings.Contains(e.CurrToken, "PATENTES") || (strings.Contains(e.CurrToken, "OTROS") && strings.Contains(e.NextToken, "OBS: N/A")){
			values := tokenize(e.CurrToken, 5)
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

/*
Function to check if a given string is or not the header of the section.
Parameter: string s
Return: True or false
*/

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

/* 
Function to load the extracted values into the OtherAsset structure. 
Parameters: values in an array of strings. The first element is not inserted because it is the index and not relevant.
Return: an instance of OtherAsset with the values from the array
*/

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

/*
Function to calculate the total of the extracted assets.
*/

func addAssets(assets []*declaration.OtherAsset) int64 {
	var total int64
	for _, a := range assets {
		total += a.Importe
	}
	return total
}
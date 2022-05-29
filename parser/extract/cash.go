package extract

import (
	"strings"
)

func Cash(e *Extractor, parser *ParserData) int64 {
	e.BindFlag(EXTRACTOR_FLAG_3)

	if e.MoveUntilContains(CurrToken, "1. ACTIVOS"){
		for e.Scan(){
			if strings.Contains(e.CurrToken, "1.1 EFECTIVO EN GS."){
				if isNumber(e.NextToken) {
					return StringToInt64(e.NextToken)
				}
				
			}
			
		}
	}

	return 0
	
}

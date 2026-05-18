package translateSchema

import "MyProject/models/translator/dataModel"

type DetailTranslator struct {
	Massage string               `json:"massage"`
	Detail  dataModel.Translator `json:"detail"`
}

type ListTranslator struct {
	Massage string                 `json:"massage"`
	List    []dataModel.Translator `json:"list"`
	Total   int                    `json:"total"`
}

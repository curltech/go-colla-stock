package service

//
//func LoadModel(modelPath string, modelNames []string) *tf.SavedModel {
//	model, err := tf.LoadSavedModel(modelPath, modelNames, nil) // 载入模型
//	if err != nil {
//		log.Fatal("LoadSavedModel(): %v", err)
//	}
//
//	log.Println("List possible ops in graphs") // 打印出所有的Operator
//	for _, op := range model.Graph.Operations() {
//		//log.Printf("Op name: %v, on device: %v", op.Name(), op.Device())
//		log.Printf("Op name: %v", op.Name())
//	}
//	return model
//}
//
//func Run() {
//	//m := LoadModel("../freeze_model", []string{"serve"})
//	//s := m.Session
//	//var json map[string]int64
//	//ret, err := s.Run(MapGraphInputs(CreateMapFromJSON(json), m),
//	//	GetGraphOutputs([]string{"prob"}, m), nil)
//	//if err != nil {
//	//	log.Fatal("Error in executing graph...", err)
//	//}
//	// ...
//}

actions
submitSp
	| response |
	
	self fillModelWithFormContent.
	
	self model: workingModel.
	self announcingObject valueChanged.
	

	response := (ZnClient new)
   url: 'http://localhost:8081/superpixels';
   entity: (ZnEntity
               with: '{"InputPath": "' , self model spInputPath , '","Type": "' , self model spParallel asString , '"}'
               type: ZnMimeType applicationJson );
   post;
   response.
	
	self psoInputPathInput text: response contents allButLast allButLast allButFirst .
actions
submitPso
	| response |
	
	self fillModelWithFormContent.
	
	self model: workingModel.
	self announcingObject valueChanged.
	
	response := (ZnClient new)
   url: 'http://localhost:8082/pso';
   entity: (ZnEntity
               with: '{"InputPath": "' , self model psoInputPath , 
							'","OutputPath": "' , self model psoOutputPath , 
							'","Type": "' , self model psoParallel asString , 
							'","NumParticles": ' , self model numParticles asString , 
							',"MaxIter": ' , self model maxIter asString ,
							',"NumThresholds": ' , self model numThresholds asString , '}'
               type: ZnMimeType applicationJson );
   post;
   response.

	self psoPosLocation: response contents allButLast allButLast allButFirst .
	
	(self model numThresholds == 1)
		ifTrue: [ self psoShowResults1Threshold ] .
	(self model numThresholds == 2)
		ifTrue: [ self psoShowResults2Thresholds ]	.
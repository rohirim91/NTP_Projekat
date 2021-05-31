model updates
fillFormWithWorkingModel

	| aModel |
	
	aModel := self workingModel.
	
	self spInputPathInput text: aModel spInputPath .
	self spParallelCheckboxInput state: aModel spParallel.
	self psoInputPathInput text: aModel psoInputPath .
	self psoParallelCheckboxInput state: aModel psoParallel .
	self numParticlesInput text: aModel numParticles asString .
	self maxIterInput text: aModel maxIter asString .
	self numThresholdsInput text: aModel numThresholds asString .
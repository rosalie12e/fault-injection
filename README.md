Fault-Injection is a module for injecting faults into Search-NDC components for Team Sourcing. 

Current faults available include:
- Latency (TFM_5001/5002: Context Deadline Exceeded) - forces an API timeout by sleeping for longer than the WebServiceTimeout. Called in Search() function after timeStart is called for 3P response time. 

How to install:
1. Add FaultInjectionParams to config.go (FaultInjectionParams map[string]interface{} `json:"FAULT_INJECTION_PARAM,omitempty"`).
2. Install fault-injection module using 'go get'.
3. Call InjectFault method within the connector code, passing the type of fault, any values to be manipulated, and the connector requestConfig. Eg for latency: "fault.InjectFault(fault.Latency, nil, requestConfig)".
4. Add FaultInjectionParams to the parameter store. 

How to use:
1. Set IsEnabled to true in the parameter store.
2. Set FailureType to desired fault type. 
3. Run tests as normal. 
4. Desired fault should trigger.

How to add new fault type:
1. In config.go, add the faulttype to the const list: FaultName = "faultname"
2. In faultfunctions.go, define a concrete implementation for the fault type: type <faultname> struct{}
3. In faultfunctions.go, define a method Execute for the concrete fault implementation. Include the fault implementation logic within this method. Pass any error back to InjectFault rather than handling it within the Execute method. 
4. In faultfunctions.go, in faultFactory, add the switch case for the new fault, mapping it to the concrete implementation.

There is a demo fault *demoFault* in the module that can be used as a template for new faults.

Notes:
Use ConventialCommit format for commit messages https://www.conventionalcommits.org/en/v1.0.0/
Uses commit-and-tag-version for release management. Run 'npm run release' to release new version.
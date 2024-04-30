Fault-Injection is a module for injecting faults into Search-NDC components for Team Sourcing. 

Current faults available include:
- Latency (TFM_5001/5002: Context Deadline Exceeded) - forces an API timeout by sleeping for longer than the WebServiceTimeout. Called in Search() function after timeStart is called for 3P response time. 

How to install:
1. Add FaultInjectionParams to config.go (FaultInjectionParams map[string]interface{} `json:"FAULT_INJECTION_PARAM,omitempty"`).
2. Install fault-module using 'go get'.
3. Call InjectFault method within the connector code, passing the type of fault, any values to be manipulated, and the connector requestConfig. Eg for latency: "fault.InjectFault(fault.Latency, nil, requestConfig)".
4. Add FaultInjectionParams to the parameter store. 

How to use:
1. Set IsEnabled to true in the parameter store.
2. Set FailureType to desired fault type. 
3. Run tests as normal. 
4. Desired fault should trigger.

Notes:
Use ConventialCommit format for commit messages https://www.conventionalcommits.org/en/v1.0.0/
Uses commit-and-tag-version for release management. Run 'npm run release' to release new version.
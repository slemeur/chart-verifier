# Propagate User Defined Configuration To Checks

* Proposal N.: 99999
* Authors: isuttonl@redhat.com
* Status: **Draft**

## Abstract

There are multiple use cases where propagating user defined configuration to checks is useful: to specify an OpenShift
version to verify the chart compatibility or influence the severity the Helm linter check should consider a failure; in
other words, any occasion the user needs to parametrize a check.

## Rationale

`chart-verifier` is the command (either through a Container Runtime or directly using the program directly) users
interface to provide information regarding the verifications to be performed for a given chart, so it is expected that
different parameters can be used at different moments in time.

Those parameters can be given to the program in two ways: through the configuration file, which is already supported by
not yet used; and through command line flags that can be used to overwrite a value defined by the configuration file.

Having both mechanisms to influence the verification session is very useful for a couple of reasons:

1. Configuration files can be distributed as *profiles*, for example one for OpenShift 4.6, another for OpenShift 4.7;
1. Profile defaults can be overridden, helping developers and other power users to debug, change checks parameters
   before committing to write a configuration/profile file.
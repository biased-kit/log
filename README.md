# The Log
The log serves to output some info. The info is related to different aspect and serves various goals. Let's consider for example following info types:
 * Errors - this info says that end user wasn't lucky to get result. it should include the details helping to investigate the problem.
 * Warning - says that the systems is not working as expected, and it has to use backup plan to return the right result. This is the errors that system can handle.
 * Info - sometimes you app should provide kind of "progress bar" info. Especially if you develop a console app. Usually you want to print this to stdout while other types are printed to stderr.
 * Debug - this messages should be still well formated and readable. Debug=Verbose. Keep in mind that there will be an other person who read it and try to solve a problem. 

 Thought it has similar names with classical log levels, it is not the same.
 Each type has own input params and behavior. The Debug info is not outputed by default, one need to put special context (WithDebug) to make it happen.
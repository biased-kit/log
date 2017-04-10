# The Log
This package provides helpers to conveniently put errors and debug info into tracer instance (see https://github.com/opentracing/opentracing-go). 

## Log errors
An error should be of type biased-kit/errors.E, it include error context and stack trace.
An error should be recorded within the span (see opentracing spec) that raises the error. I.e. when you're going to transmit an error to caller goroutine don't forget to execute logger.

## Log debug
By default Debug recording is off and could be turned on by context (see WithDebug func).


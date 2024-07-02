package try

// RethrowPanic Special re-throw panic
const RethrowPanic = "___throw_it___"

// Define types of Try-Finally-Catch
type (
	// F Function type
	F func()
	// E Error type
	E interface{}
	// EF Error function type
	EF func(err E)
	// It structure
	It struct {
		finally F
		Error   E
	}
)

// Throw function (return or rethrow an exception)
func Throw(e E) {
	// Throw default error
	if e == nil {
		panic(RethrowPanic)
	} else {
		// Throw a specific exception
		panic(e)
	}
}

// Perform register the main-logic function.
func Perform(funcToTry F) (o *It) {
	// Initial exception object with null values
	o = &It{nil, nil}

	// Catch throw in from main logic
	defer func() {
		o.Error = recover()
	}()

	// Perform main logic
	funcToTry()

	// Response instance of It instance
	return
}

// Finally register the finally-logic function.
func (o *It) Finally(finallyFunc F) *It {
	if o.finally != nil {
		panic("Finally Function by default !!")
	} else {
		o.finally = finallyFunc
	}

	return o
}

// Catch register the finally-logic function.
func (o *It) Catch(funcCaught EF) *It {
	// Check if it has Error
	if o.Error != nil {
		// Catch error in from catching logic
		defer func() {
			// Call finally (Before receive error from recovering process)
			if o.finally != nil {
				o.finally()
			}

			// Receive error from recovering process
			if err := recover(); err != nil {
				// If it is just re-throw panic Exception
				if err == RethrowPanic {
					err = o.Error
				}
				panic(err)
			}
		}()

		// Perform catching logic
		funcCaught(o.Error)
	} else if o.finally != nil {
		// Perform finally logic
		o.finally()
	}

	return o
}

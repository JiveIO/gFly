package try

import (
	"testing"
)

func Test_NormalFlow(t *testing.T) {
	called := false

	Perform(func() {
		called = true
	}).Catch(func(_ E) {
		t.Error("Catch must not be called")
	})

	// if try was not called
	if !called {
		t.Error("Try do not called")
	}
}

func Test_NormalFlowFinally(t *testing.T) {
	calledTry := false
	calledFinally := false

	Perform(func() {
		calledTry = true
	}).Finally(func() {
		calledFinally = true
	}).Catch(func(_ E) {
		t.Error("Catch must not be called")
	})

	// if try was not called
	if !calledTry {
		t.Error("Try do not called")
	}

	// if finally was not called
	if !calledFinally {
		t.Error("Finally do not called")
	}
}

func Test_CrashInTry(t *testing.T) {
	calledFinally := false
	calledCatch := false

	Perform(func() {
		panic("testing panic")
	}).Finally(func() {
		calledFinally = true
	}).Catch(func(e E) {
		calledCatch = true
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
	})

	// if catch was not called
	if !calledCatch {
		t.Error("Catch do not called")
	}

	// if finally was not called
	if !calledFinally {
		t.Error("Finally do not called")
	}
}

func Test_CrashInTry2(t *testing.T) {
	calledCatch := false

	Perform(func() {
		panic("testing panic")
	}).Catch(func(e E) {
		calledCatch = true
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
	})

	// if catch was not called
	if !calledCatch {
		t.Error("Catch do not called")
	}
}

func Test_CrashInCatch(t *testing.T) {
	calledFinally := false

	defer func() {
		err := recover()

		if err != "another panic" {
			t.Error("error is not 'another panic'")
		}
		// if finally was not called
		if !calledFinally {
			t.Error("Finally do not called")
		}
	}()

	Perform(func() {
		panic("testing panic")
	}).Finally(func() {
		calledFinally = true
	}).Catch(func(e E) {
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
		panic("another panic")
	})
}

func Test_CrashInCatch2(t *testing.T) {
	defer func() {
		err := recover()
		if err != "another panic" {
			t.Error("error is not 'another panic'")
		}
	}()
	Perform(func() {
		panic("testing panic")

	}).Catch(func(e E) {
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
		panic("another panic")
	})
}

func Test_CrashInThrow(t *testing.T) {
	calledFinally := false

	defer func() {
		err := recover()
		if err != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
		// if finally was not called
		if !calledFinally {
			t.Error("Finally do not called")
		}
	}()

	Perform(func() {
		panic("testing panic")
	}).Finally(func() {
		calledFinally = true
	}).Catch(func(e E) {
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
		Throw(nil)
	})
}

func Test_CrashInThrow2(t *testing.T) {
	defer func() {
		err := recover()
		if err != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
	}()

	Perform(func() {
		panic("testing panic")
	}).Catch(func(e E) {
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
		Throw(nil)
	})
}

func Test_CrashInFinally1(t *testing.T) {
	calledTry := false

	defer func() {
		err := recover()
		if err != "finally panic" {
			t.Error("error is not 'finally panic'")
		}

		// if try was not called
		if !calledTry {
			t.Error("Try do not called")
		}
	}()

	Perform(func() {
		calledTry = true
	}).Finally(func() {
		panic("finally panic")
	}).Catch(func(e E) {
		t.Error("Catch must not be called")
	})
}

func Test_CrashInFinally2(t *testing.T) {
	defer func() {
		err := recover()
		if err != "finally panic" {
			t.Error("error is not 'finally panic'")
		}
	}()

	Perform(func() {
		panic("testing panic")
	}).Finally(func() {
		panic("finally panic")
	}).Catch(func(e E) {
		if e != "testing panic" {
			t.Error("error is not 'testing panic'")
		}
		panic("another panic")
	})
}

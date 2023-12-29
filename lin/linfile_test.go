package lin_test

import (
	"acdc/lin"
	"testing"
)

func TestReadLinFile(t *testing.T) {

	ld, err := lin.ReadLinFile("testdata/StC_test_OC4Semi_Linear_Tow.1.lin")
	if err != nil {
		t.Fatal(err)
	}

	if act, exp := ld.SimTime, 0.0; act != exp {
		t.Fatalf("ld.SimTime = %v, expected %v", act, exp)
	}
	if act, exp := ld.RotorSpeed, 1.2566; act != exp {
		t.Fatalf("ld.RotorSpeed = %v, expected %v", act, exp)
	}
	if act, exp := ld.Azimuth, 0.0; act != exp {
		t.Fatalf("ld.Azimuth = %v, expected %v", act, exp)
	}
	if act, exp := ld.WindSpeed, 0.0; act != exp {
		t.Fatalf("ld.WindSpeed = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_x, 8; act != exp {
		t.Fatalf("ld.Num_x = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_xd, 0; act != exp {
		t.Fatalf("ld.Num_xd = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_z, 0; act != exp {
		t.Fatalf("ld.Num_z = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_u, 6; act != exp {
		t.Fatalf("ld.Num_u = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_y, 19; act != exp {
		t.Fatalf("ld.Num_y = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_x), 8; act != exp {
		t.Fatalf("len(ld.OP_x) = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_xdot), 8; act != exp {
		t.Fatalf("len(ld.OP_xdot) = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_u), 6; act != exp {
		t.Fatalf("len(ld.OP_u) = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_y), 19; act != exp {
		t.Fatalf("len(ld.OP_y) = %v, expected %v", act, exp)
	}
	rAct, cAct := ld.A.Dims()
	if rExp, cExp := 8, 8; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.A.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
	rAct, cAct = ld.B.Dims()
	if rExp, cExp := 8, 6; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.B.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
	rAct, cAct = ld.C.Dims()
	if rExp, cExp := 19, 8; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.C.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
	rAct, cAct = ld.D.Dims()
	if rExp, cExp := 19, 6; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.D.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}

	ld, err = lin.ReadLinFile("testdata/Ideal_Beam_Fixed_Free_Linear.1.lin")
	if err != nil {
		t.Fatal(err)
	}

	if act, exp := ld.SimTime, 0.0; act != exp {
		t.Fatalf("ld.SimTime = %v, expected %v", act, exp)
	}
	if act, exp := ld.RotorSpeed, 0.0; act != exp {
		t.Fatalf("ld.RotorSpeed = %v, expected %v", act, exp)
	}
	if act, exp := ld.Azimuth, 0.0; act != exp {
		t.Fatalf("ld.Azimuth = %v, expected %v", act, exp)
	}
	if act, exp := ld.WindSpeed, 0.0; act != exp {
		t.Fatalf("ld.WindSpeed = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_x, 96; act != exp {
		t.Fatalf("ld.Num_x = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_xd, 0; act != exp {
		t.Fatalf("ld.Num_xd = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_z, 0; act != exp {
		t.Fatalf("ld.Num_z = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_u, 166; act != exp {
		t.Fatalf("ld.Num_u = %v, expected %v", act, exp)
	}
	if act, exp := ld.Num_y, 327; act != exp {
		t.Fatalf("ld.Num_y = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_x), 96; act != exp {
		t.Fatalf("len(ld.OP_x) = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_xdot), 96; act != exp {
		t.Fatalf("len(ld.OP_xdot) = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_u), 166; act != exp {
		t.Fatalf("len(ld.OP_u) = %v, expected %v", act, exp)
	}
	if act, exp := len(ld.OP_y), 327; act != exp {
		t.Fatalf("len(ld.OP_y) = %v, expected %v", act, exp)
	}
	rAct, cAct = ld.A.Dims()
	if rExp, cExp := 96, 96; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.A.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
	rAct, cAct = ld.B.Dims()
	if rExp, cExp := 96, 166; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.B.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
	rAct, cAct = ld.C.Dims()
	if rExp, cExp := 327, 96; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.C.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
	rAct, cAct = ld.D.Dims()
	if rExp, cExp := 327, 166; rAct != rExp || cAct != cExp {
		t.Fatalf("ld.D.Dims() = [%v,%v], expected [%v,%v]", rAct, cAct, rExp, cExp)
	}
}

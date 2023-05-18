package fio

import (
	"acdc/fio/schema"
	"bytes"
	"fmt"
)

type File struct {
	Path
	AeroDyn14      *AeroDyn14
	AeroDyn15      *AeroDyn15
	AeroDynBlade   *AeroDynBlade
	AirfoilInfo    *AirfoilInfo
	BeamDyn        *BeamDyn
	BeamDynBlade   *BeamDynBlade
	ElastoDyn      *ElastoDyn
	ElastoDynBlade *ElastoDynBlade
	ElastoDynTower *ElastoDynTower
	FreeVortexWake *FreeVortexWake
	InflowWind     *InflowWind
	Main           *Main
	ServoDyn       *ServoDyn
	SubDyn         *SubDyn
	TailFin        *TailFin
	TextFile       *TextFile
	UniformWind    *UniformWind
}

type Files struct {
	Paths
	AeroDyn14      []*AeroDyn14
	AeroDyn15      []*AeroDyn15
	AeroDynBlade   []*AeroDynBlade
	AirfoilInfo    []*AirfoilInfo
	BeamDyn        []*BeamDyn
	BeamDynBlade   []*BeamDynBlade
	ElastoDyn      []*ElastoDyn
	ElastoDynBlade []*ElastoDynBlade
	ElastoDynTower []*ElastoDynTower
	FreeVortexWake []*FreeVortexWake
	InflowWind     []*InflowWind
	Main           []*Main
	ServoDyn       []*ServoDyn
	SubDyn         []*SubDyn
	TailFin        []*TailFin
	TextFile       []*TextFile
	UniformWind    []*UniformWind
}

//------------------------------------------------------------------------------
// AeroDyn14
//------------------------------------------------------------------------------

type AeroDyn14 struct {
	Header1      Header           // AeroDyn14 Input File
	Title        Title            //
	StallMod     String           // Dynamic stall included [BEDDOES or STEADY]
	UseCm        String           // Use aerodynamic pitching moment model? [USE_CM or NO_CM]
	InfModel     String           // Inflow model [DYNIN or EQUIL]
	IndModel     String           // Induction-factor model [NONE or WAKE or SWIRL]
	AToler       Float            // Induction-factor tolerance (convergence criteria)
	TLModel      String           // Tip-loss model (EQUIL only) [PRANDtl, GTECH, or NONE]
	HLModel      String           // Hub-loss model (EQUIL only) [PRANdtl or NONE]
	TwrShad      Float            // Tower-shadow velocity deficit
	ShadHWid     Float            // Tower-shadow half width
	T_Shad_Refpt Float            // Tower-shadow reference point
	AirDens      Float            // Air density
	KinVisc      Float            // Kinematic air viscosity [CURRENTLY IGNORED]
	DTAero       FloatDefault     // Time interval for aerodynamic calculations
	NumFoil      Int              // Number of airfoil files
	FoilNm       Files            // Names of the airfoil files [NumFoil lines]
	BldNodes     Int              // Number of blade nodes used for analysis
	BldNodeData  TableBldNodeData //
}

func (s *AeroDyn14) Parse(path string) error {
	return parse(s, path, schema.AeroDyn14)
}

func (s *AeroDyn14) Format(path string) error {
	return format(s, path, schema.AeroDyn14)
}

type TableBldNodeData struct {
	Rows []TableBldNodeDataRow
}

type TableBldNodeDataRow struct {
	RNodes   float64 //
	AeroTwst float64 //
	DRNodes  float64 //
	Chord    float64 //
	NFoil    int     //
	PrnElm   string  //
}

func (t *TableBldNodeData) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableBldNodeDataRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableBldNodeData) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// AeroDyn15
//------------------------------------------------------------------------------

type AeroDyn15 struct {
	Header1           Header       // AeroDyn15 Input File
	Title             Title        //
	Header2           Header       // General Options
	Echo              Bool         // Echo the input to "<rootname>.AD.ech"?
	DTAero            FloatDefault // Time interval for aerodynamic calculations {or "default"}
	WakeMod           Int          // Type of wake/induction model {0=none, 1=BEMT, 2=DBEMT, 3=OLAF} [WakeMod cannot be 2 or 3 when linearizing]
	AFAeroMod         Int          // Type of blade airfoil aerodynamics model {1=steady model, 2=Beddoes-Leishman unsteady model} [AFAeroMod must be 1 when linearizing]
	TwrPotent         Int          // Type tower influence on wind based on potential flow around the tower {0=none, 1=baseline potential flow, 2=potential flow with Bak correction}
	TwrShadow         Int          // Calculate tower influence on wind based on downstream tower shadow {0=none, 1=Powles model, 2=Eames model}
	TwrAero           Bool         // Calculate tower aerodynamic loads?
	FrozenWake        Bool         // Assume frozen wake during linearization? [used only when WakeMod=1 and when linearizing]
	CavitCheck        Bool         // Perform cavitation check? [AFAeroMod must be 1 when CavitCheck=true]
	Buoyancy          Bool         // Include buoyancy effects?
	CompAA            Bool         // Flag to compute AeroAcoustics calculation [used only when WakeMod = 1 or 2]
	AA_InputFile      Path         // AeroAcoustics input file [used only when CompAA=true]
	Header3           Header       // Environmental Conditions
	AirDens           FloatDefault // Air density
	KinVisc           FloatDefault // Kinematic viscosity of working fluid
	SpdSound          FloatDefault // Speed of sound in working fluid
	Patm              FloatDefault // Atmospheric pressure [used only when CavitCheck=True]
	Pvap              FloatDefault // Vapour pressure of working fluid [used only when CavitCheck=True]
	Header4           Header       // Blade-Element/Momentum Theory Options [unused when WakeMod=0 or 3]
	SkewMod           Int          // Type of skewed-wake correction model {1=uncoupled, 2=Pitt/Peters, 3=coupled} [unused when WakeMod=0 or 3]
	SkewModFactor     FloatDefault // Constant used in Pitt/Peters skewed wake model {or "default" is 15/32*pi} [used only when SkewMod=2; unused when WakeMod=0 or 3]
	TipLoss           Bool         // Use the Prandtl tip-loss model? [unused when WakeMod=0 or 3]
	HubLoss           Bool         // Use the Prandtl hub-loss model? [unused when WakeMod=0 or 3]
	TanInd            Bool         // Include tangential induction in BEMT calculations? [unused when WakeMod=0 or 3]
	AIDrag            Bool         // Include the drag term in the axial-induction calculation? [unused when WakeMod=0 or 3]
	TIDrag            Bool         // Include the drag term in the tangential-induction calculation? [unused when WakeMod=0,3 or TanInd=FALSE]
	IndToler          FloatDefault // Convergence tolerance for BEMT nonlinear solve residual equation {or "default"} [unused when WakeMod=0 or 3]
	MaxIter           Int          // Maximum number of iteration steps [unused when WakeMod=0]
	Header5           Header       // Dynamic Blade-Element/Momentum Theory Options [used only when WakeMod=2]
	DBEMT_Mod         Int          // Type of dynamic BEMT (DBEMT) model {1=constant tau1, 2=time-dependent tau1, 3=constant tau1 with continuous formulation} [used only when WakeMod=2]
	Tau1_const        Float        // Time constant for DBEMT (s) [used only when WakeMod=2 and DBEMT_Mod=1 or 3]
	Header6           Header       // OLAF -- cOnvecting LAgrangian Filaments (Free Vortex Wake) Theory Options [used only when WakeMod=3]
	OLAFInputFileName File         // Input file for OLAF [used only when WakeMod=3]
	Header7           Header       // Beddoes-Leishman Unsteady Airfoil Aerodynamics Options [used only when AFAeroMod=2]
	UAMod             Int          // Unsteady Aero Model Switch {2=B-L Gonzalez, 3=B-L Minnema/Pierce, 4=B-L HGM 4-states, 5=B-L 5 states, 6=Oye, 7=Boeing-Vertol} [used only when AFAeroMod=2]
	FLookup           Bool         // Flag to indicate whether a lookup for f' will be calculated (TRUE) or whether best-fit exponential equations will be used (FALSE); if FALSE S1-S4 must be provided in airfoil input files [used only when AFAeroMod=2]
	Header8           Header       // Airfoil Information
	AFTabMod          Int          // Interpolation method for multiple airfoil tables {1=1D interpolation on AoA (first table only); 2=2D interpolation on AoA and Re; 3=2D interpolation on AoA and UserProp}
	InCol_Alfa        Int          // The column in the airfoil tables that contains the angle of attack
	InCol_Cl          Int          // The column in the airfoil tables that contains the lift coefficient
	InCol_Cd          Int          // The column in the airfoil tables that contains the drag coefficient
	InCol_Cm          Int          // The column in the airfoil tables that contains the pitching-moment coefficient; use zero if there is no Cm column
	InCol_Cpmin       Int          // The column in the airfoil tables that contains the Cpmin coefficient; use zero if there is no Cpmin column
	NumAFfiles        Int          // Number of airfoil files used
	AFNames           Files        // Airfoil file names
	Header9           Header       // Rotor/Blade Properties
	UseBlCm           Bool         // Include aerodynamic pitching moment in calculations?
	ADBlFile1         File         // Name of file containing distributed aerodynamic properties for Blade #1
	ADBlFile2         File         // Name of file containing distributed aerodynamic properties for Blade #2 [unused if NumBl < 2]
	ADBlFile3         File         // Name of file containing distributed aerodynamic properties for Blade #3 [unused if NumBl < 3]
	Header10          Header       // Hub Properties [used only when Buoyancy=True]
	VolHub            Float        // Hub volume (m^3)
	HubCenBx          Float        // Hub center of buoyancy x direction offset (m)
	Header11          Header       // Nacelle Properties [used only when Buoyancy=True]
	VolNac            Float        // Nacelle volume
	NacCenB           Floats       // Position of nacelle center of buoyancy from yaw bearing in nacelle coordinates
	Header12          Header       // Tail fin Aerodynamics
	TFinAero          Bool         // Calculate tail fin aerodynamics model
	TFinFile          Path         // Input file for tail fin aerodynamics [used only when TFinAero=True]
	Header13          Header       // Tower Influence and Aerodynamics [used only when TwrPotent/=0, TwrShadow/=0, TwrAero=True, or Buoyancy=True]
	NumTwrNds         Int          // Number of tower nodes used in the analysis [used only when TwrPotent/=0, TwrShadow/=0, TwrAero=True, or Buoyancy=True]
	TwrNd             TableTwrNd   //
	Header14          Header       // Outputs
	SumPrint          Bool         // Generate a summary file listing input options and interpolated properties to "<rootname>.AD.sum"?
	NBlOuts           Int          // Number of blade node outputs [0 - 9]
	BlOutNd           Ints         // Blade nodes whose values will be output
	NTwOuts           Int          // Number of tower node outputs [0 - 9]
	TwOutNd           Ints         // Tower nodes whose values will be output
	OutList           OutList      // The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels
}

func (s *AeroDyn15) Parse(path string) error {
	return parse(s, path, schema.AeroDyn15)
}

func (s *AeroDyn15) Format(path string) error {
	return format(s, path, schema.AeroDyn15)
}

type TableTwrNd struct {
	Rows []TableTwrNdRow
}

type TableTwrNdRow struct {
	TwrElev float64 // m
	TwrDiam float64 // m
	TwrCd   float64 // -
	TwrTI   float64 // -
	TwrCb   float64 // -
}

func (t *TableTwrNd) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableTwrNdRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableTwrNd) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// AeroDynBlade
//------------------------------------------------------------------------------

type AeroDynBlade struct {
	Header1  Header     // AeroDynBlade Input File
	Title    Title      //
	Header2  Header     // Blade Properties
	NumBlNds Int        // Number of blade nodes used in the analysis
	BlNds    TableBlNds //
}

func (s *AeroDynBlade) Parse(path string) error {
	return parse(s, path, schema.AeroDynBlade)
}

func (s *AeroDynBlade) Format(path string) error {
	return format(s, path, schema.AeroDynBlade)
}

type TableBlNds struct {
	Rows []TableBlNdsRow
}

type TableBlNdsRow struct {
	BlSpn    float64 // m
	BlCrvAC  float64 // m
	BlSwpAC  float64 // m
	BlCrvAng float64 // deg
	BlTwist  float64 // deg
	BlChord  float64 // m
	BlAFID   int     // -
	BlCb     float64 // -
	BlCenBn  float64 // m
	BlCenBt  float64 // m
}

func (t *TableBlNds) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableBlNdsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableBlNds) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// AirfoilInfo
//------------------------------------------------------------------------------

type AirfoilInfo struct {
	Text Text //
}

func (s *AirfoilInfo) Parse(path string) error {
	return parse(s, path, schema.AirfoilInfo)
}

func (s *AirfoilInfo) Format(path string) error {
	return format(s, path, schema.AirfoilInfo)
}

//------------------------------------------------------------------------------
// BeamDyn
//------------------------------------------------------------------------------

type BeamDyn struct {
	Header1          Header        // BeamDyn Input File
	Title            Title         //
	Header2          Header        // Simulation Control
	Echo             Bool          // Echo input data to "<RootName>.ech"?
	QuasiStaticInit  Bool          // Use quasi-static pre-conditioning with centripetal accelerations in initialization? [dynamic solve only]
	Rhoinf           Float         // Numerical damping parameter for generalized-alpha integrator
	Quadrature       Int           // Quadrature method: 1=Gaussian; 2=Trapezoidal
	Refine           IntDefault    // Refinement factor for trapezoidal quadrature [DEFAULT = 1; used only when quadrature=2]
	N_fact           IntDefault    // Factorization frequency for the Jacobian in N-R iteration [DEFAULT = 5]
	DTBeam           FloatDefault  // Time step size
	Load_retries     IntDefault    // Number of factored load retries before quitting the simulation [DEFAULT = 20]
	NRMax            IntDefault    // Max number of iterations in Newton-Raphson algorithm [DEFAULT = 10]
	Stop_tol         FloatDefault  // Tolerance for stopping criterion [DEFAULT = 1E-5]
	Tngt_stf_fd      BoolDefault   // Use finite differenced tangent stiffness matrix?
	Tngt_stf_comp    BoolDefault   // Compare analytical finite differenced tangent stiffness matrix?
	Tngt_stf_pert    FloatDefault  // Perturbation size for finite differencing [DEFAULT = 1E-6]
	Tngt_stf_difftol FloatDefault  // Maximum allowable relative difference between analytical and fd tangent stiffness; [DEFAULT = 0.1]
	RotStates        Bool          // Orient states in the rotating frame during linearization? [used only when linearizing]
	Header3          Header        // Geometry Parameter
	Member_total     Int           // Total number of members
	Kp_total         Int           // Total number of key points [must be at least 3]
	KPMember         TableKPMember //
	KPData           TableKPData   //
	Header4          Header        // Mesh Parameter
	Order_elem       Int           // Order of interpolation (basis) function
	Header5          Header        // Material Parameter
	BldFile          File          // Name of file containing properties for blade
	Header6          Header        // Pitch Actuator Parameters
	UsePitchAct      Bool          // Whether a pitch actuator should be used
	PitchJ           Float         // Pitch actuator inertia [used only when UsePitchAct is true]
	PitchK           Float         // Pitch actuator stiffness [used only when UsePitchAct is true]
	PitchC           Float         // Pitch actuator damping [used only when UsePitchAct is true]
	Header7          Header        // Outputs
	SumPrint         Bool          // Print summary data to "<RootName>.sum"
	OutFmt           String        // Format used for text tabular output (except time)
	NNodeOuts        Int           // Number of nodes to output to file [0 - 9]
	OutNd            Ints          // Nodes whose values will be output
	OutList          OutList       // The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels
}

func (s *BeamDyn) Parse(path string) error {
	return parse(s, path, schema.BeamDyn)
}

func (s *BeamDyn) Format(path string) error {
	return format(s, path, schema.BeamDyn)
}

type TableKPMember struct {
	Rows []TableKPMemberRow
}

type TableKPMemberRow struct {
	Member int // -
	NumKP  int // -
}

func (t *TableKPMember) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableKPMemberRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableKPMember) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableKPData struct {
	Rows []TableKPDataRow
}

type TableKPDataRow struct {
	Kp_xr         float64 // m
	Kp_yr         float64 // m
	Kp_zr         float64 // m
	Initial_twist float64 // deg
}

func (t *TableKPData) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableKPDataRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableKPData) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// BeamDynBlade
//------------------------------------------------------------------------------

type BeamDynBlade struct {
	Header1       Header     // BeamDynBlade Input File
	Title         Title      //
	Header2       Header     // Blade Parameters
	Station_total Int        // Number of blade input stations (-)
	Damp_type     Int        // Damping type: 0: no damping; 1: damped
	Header3       Header     // Damping Coefficient
	Mu            TableMu    //
	Header4       Header     // Distributed Properties
	Stations      BDStations //
}

func (s *BeamDynBlade) Parse(path string) error {
	return parse(s, path, schema.BeamDynBlade)
}

func (s *BeamDynBlade) Format(path string) error {
	return format(s, path, schema.BeamDynBlade)
}

type TableMu struct {
	Rows []TableMuRow
}

type TableMuRow struct {
	Mu1 float64 // -
	Mu2 float64 // -
	Mu3 float64 // -
	Mu4 float64 // -
	Mu5 float64 // -
	Mu6 float64 // -
}

func (t *TableMu) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableMuRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableMu) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// ElastoDyn
//------------------------------------------------------------------------------

type ElastoDyn struct {
	Header1   Header       // ElastoDyn Input File
	Title     Title        //
	Header2   Header       // Simulation Control
	Echo      Bool         // Echo input data to <RootName>.ech
	Method    Int          // Integration method: {1: RK4, 2: AB4, or 3: ABM4}
	DT        FloatDefault // Integration time step
	Header3   Header       // Degrees of Freedom
	FlapDOF1  Bool         // First flapwise blade mode DOF
	FlapDOF2  Bool         // Second flapwise blade mode DOF
	EdgeDOF   Bool         // First edgewise blade mode DOF
	TeetDOF   Bool         // Rotor-teeter DOF (flag) [unused for 3 blades]
	DrTrDOF   Bool         // Drivetrain rotational-flexibility DOF
	GenDOF    Bool         // Generator DOF
	YawDOF    Bool         // Yaw DOF
	TwFADOF1  Bool         // First fore-aft tower bending-mode DOF
	TwFADOF2  Bool         // Second fore-aft tower bending-mode DOF
	TwSSDOF1  Bool         // First side-to-side tower bending-mode DOF
	TwSSDOF2  Bool         // Second side-to-side tower bending-mode DOF
	PtfmSgDOF Bool         // Platform horizontal surge translation DOF
	PtfmSwDOF Bool         // Platform horizontal sway translation DOF
	PtfmHvDOF Bool         // Platform vertical heave translation DOF
	PtfmRDOF  Bool         // Platform roll tilt rotation DOF
	PtfmPDOF  Bool         // Platform pitch tilt rotation DOF
	PtfmYDOF  Bool         // Platform yaw rotation DOF
	Header4   Header       // Initial Conditions
	OoPDefl   Float        // Initial out-of-plane blade-tip displacement
	IPDefl    Float        // Initial in-plane blade-tip deflection
	BlPitch1  Float        // Blade 1 initial pitch
	BlPitch2  Float        // Blade 2 initial pitch
	BlPitch3  Float        // Blade 3 initial pitch [unused for 2 blades]
	TeetDefl  Float        // Initial or fixed teeter angle [unused for 3 blades]
	Azimuth   Float        // Initial azimuth angle for blade 1
	RotSpeed  Float        // Initial or fixed rotor speed
	NacYaw    Float        // Initial or fixed nacelle-yaw angle
	TTDspFA   Float        // Initial fore-aft tower-top displacement
	TTDspSS   Float        // Initial side-to-side tower-top displacement
	PtfmSurge Float        // Initial or fixed horizontal surge translational displacement of platform
	PtfmSway  Float        // Initial or fixed horizontal sway translational displacement of platform
	PtfmHeave Float        // Initial or fixed vertical heave translational displacement of platform
	PtfmRoll  Float        // Initial or fixed roll tilt rotational displacement of platform
	PtfmPitch Float        // Initial or fixed pitch tilt rotational displacement of platform
	PtfmYaw   Float        // Initial or fixed yaw rotational displacement of platform
	Header5   Header       // Turbine Configuration
	NumBl     Int          // Number of blades
	TipRad    Float        // The distance from the rotor apex to the blade tip
	HubRad    Float        // The distance from the rotor apex to the blade root
	PreCone1  Float        // Blade 1 cone angle
	PreCone2  Float        // Blade 2 cone angle
	PreCone3  Float        // Blade 3 cone angle [unused for 2 blades]
	HubCM     Float        // Distance from rotor apex to hub mass [positive downwind]
	UndSling  Float        // Undersling length [distance from teeter pin to the rotor apex] (meters) [unused for 3 blades]
	Delta3    Float        // Delta-3 angle for teetering rotors [unused for 3 blades]
	AzimB1Up  Float        // Azimuth value to use for I/O when blade 1 points up
	OverHang  Float        // Distance from yaw axis to rotor apex [3 blades] or teeter pin [2 blades]
	ShftGagL  Float        // Distance from rotor apex [3 blades] or teeter pin [2 blades] to shaft strain gages [positive for upwind rotors]
	ShftTilt  Float        // Rotor shaft tilt angle
	NacCMxn   Float        // Downwind distance from the tower-top to the nacelle CM
	NacCMyn   Float        // Lateral  distance from the tower-top to the nacelle CM
	NacCMzn   Float        // Vertical distance from the tower-top to the nacelle CM
	NcIMUxn   Float        // Downwind distance from the tower-top to the nacelle IMU
	NcIMUyn   Float        // Lateral  distance from the tower-top to the nacelle IMU
	NcIMUzn   Float        // Vertical distance from the tower-top to the nacelle IMU
	Twr2Shft  Float        // Vertical distance from the tower-top to the rotor shaft
	TowerHt   Float        // Height of tower above ground level [onshore], MSL [offshore], or seabed [MHK]
	TowerBsHt Float        // Height of tower base above ground level [onshore], MSL [offshore], or seabed [MHK]
	PtfmCMxt  Float        // Downwind distance from the ground level [onshore], MSL [offshore], or seabed [MHK] to the platform CM
	PtfmCMyt  Float        // Lateral distance from the ground level [onshore], MSL [offshore], or seabed [MHK] to the platform CM
	PtfmCMzt  Float        // Vertical distance from the ground level [onshore], MSL [offshore], or seabed [MHK] to the platform CM
	PtfmRefzt Float        // Vertical distance from the ground level [onshore], MSL [offshore], or seabed [MHK] to the platform reference point
	Header6   Header       // Mass and Inertia
	TipMass1  Float        // Tip-brake mass, blade 1
	TipMass2  Float        // Tip-brake mass, blade 2
	TipMass3  Float        // Tip-brake mass, blade 3 [unused for 2 blades]
	HubMass   Float        // Hub mass
	HubIner   Float        // Hub inertia about rotor axis [3 blades] or teeter axis [2 blades]
	GenIner   Float        // Generator inertia about HSS
	NacMass   Float        // Nacelle mass
	NacYIner  Float        // Nacelle inertia about yaw axis
	YawBrMass Float        // Yaw bearing mass
	PtfmMass  Float        // Platform mass
	PtfmRIner Float        // Platform inertia for roll tilt rotation about the platform CM
	PtfmPIner Float        // Platform inertia for pitch tilt rotation about the platform CM
	PtfmYIner Float        // Platform inertia for yaw rotation about the platform CM
	Header7   Header       // Blade
	BldNodes  Int          // Number of blade nodes (per blade) used for analysis
	BldFile1  File         // Name of file containing properties for blade 1
	BldFile2  File         // Name of file containing properties for blade 2
	BldFile3  File         // Name of file containing properties for blade 3 [unused for 2 blades]
	Header8   Header       // Rotor-Teeter
	TeetMod   Int          // Rotor-teeter spring/damper model {0: none, 1: standard, 2: user-defined from routine UserTeet} [unused for 3 blades]
	TeetDmpP  Float        // Rotor-teeter damper position [used only for 2 blades and when TeetMod=1]
	TeetDmp   Float        // Rotor-teeter damping constant [used only for 2 blades and when TeetMod=1]
	TeetCDmp  Float        // Rotor-teeter rate-independent Coulomb-damping moment [used only for 2 blades and when TeetMod=1]
	TeetSStP  Float        // Rotor-teeter soft-stop position [used only for 2 blades and when TeetMod=1]
	TeetHStP  Float        // Rotor-teeter hard-stop position [used only for 2 blades and when TeetMod=1]
	TeetSSSp  Float        // Rotor-teeter soft-stop linear-spring constant [used only for 2 blades and when TeetMod=1]
	TeetHSSp  Float        // Rotor-teeter hard-stop linear-spring constant [used only for 2 blades and when TeetMod=1]
	Header9   Header       // Drivetrain
	GBoxEff   Float        // Gearbox efficiency
	GBRatio   Float        // Gearbox ratio
	DTTorSpr  Float        // Drivetrain torsional spring
	DTTorDmp  Float        // Drivetrain torsional damper
	Header10  Header       // Furling
	Furling   Bool         // Read in additional model properties for furling turbine [must currently be FALSE)
	FurlFile  Path         // Name of file containing furling properties [unused when Furling=False]
	Header11  Header       // Tower
	TwrNodes  Int          // Number of tower nodes used for analysis
	TwrFile   File         // Name of file containing tower properties
	Header12  Header       // Output
	SumPrint  Bool         // Print summary data to "<RootName>.sum"
	OutFile   Int          // Switch to determine where output will be placed: {1: in module output file only; 2: in glue code output file only; 3: both} (currently unused)
	TabDelim  Bool         // Use tab delimiters in text tabular output file? (currently unused)
	OutFmt    String       // Format used for text tabular output (except time).  Resulting field should be 10 characters. (currently unused)
	TStart    Float        // Time to begin tabular output (currently unused)
	DecFact   Int          // Decimation factor for tabular output {1: output every time step} (currently unused)
	NTwGages  Int          // Number of tower nodes that have strain gages for output [0 to 9]
	TwrGagNd  Ints         // List of tower nodes that have strain gages [1 to TwrNodes] [unused if NTwGages=0]
	NBlGages  Int          // Number of blade nodes that have strain gages for output [0 to 9]
	BldGagNd  Ints         // List of blade nodes that have strain gages [1 to BldNodes] [unused if NBlGages=0]
	OutList   OutList      // The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels
}

func (s *ElastoDyn) Parse(path string) error {
	return parse(s, path, schema.ElastoDyn)
}

func (s *ElastoDyn) Format(path string) error {
	return format(s, path, schema.ElastoDyn)
}

//------------------------------------------------------------------------------
// ElastoDynBlade
//------------------------------------------------------------------------------

type ElastoDynBlade struct {
	Header1   Header       // ElastoDynBlade Input File
	Title     Title        //
	Header2   Header       // Blade Parameters
	NBlInpSt  Int          // Number of blade input stations
	BldFlDmp1 Float        // Blade flap mode #1 structural damping in percent of critical
	BldFlDmp2 Float        // Blade flap mode #2 structural damping in percent of critical
	BldEdDmp1 Float        // Blade edge mode #1 structural damping in percent of critical
	Header3   Header       // Blade Adjustment Factors
	FlStTunr1 Float        // Blade flapwise modal stiffness tuner, 1st mode
	FlStTunr2 Float        // Blade flapwise modal stiffness tuner, 2nd mode
	AdjBlMs   Float        // Factor to adjust blade mass density
	AdjFlSt   Float        // Factor to adjust blade flap stiffness
	AdjEdSt   Float        // Factor to adjust blade edge stiffness
	Header4   Header       // Distributed Blade Properties
	BlInpSt   TableBlInpSt //
	Header5   Header       // Blade Mode Shapes
	BldFl1Sh2 Float        // Flap mode 1, coeff of x^2
	BldFl1Sh3 Float        //            , coeff of x^3
	BldFl1Sh4 Float        //            , coeff of x^4
	BldFl1Sh5 Float        //            , coeff of x^5
	BldFl1Sh6 Float        //            , coeff of x^6
	BldFl2Sh2 Float        // Flap mode 2, coeff of x^2
	BldFl2Sh3 Float        //            , coeff of x^3
	BldFl2Sh4 Float        //            , coeff of x^4
	BldFl2Sh5 Float        //            , coeff of x^5
	BldFl2Sh6 Float        //            , coeff of x^6
	BldEdgSh2 Float        // Edge mode 1, coeff of x^2
	BldEdgSh3 Float        //            , coeff of x^3
	BldEdgSh4 Float        //            , coeff of x^4
	BldEdgSh5 Float        //            , coeff of x^5
	BldEdgSh6 Float        //            , coeff of x^6
}

func (s *ElastoDynBlade) Parse(path string) error {
	return parse(s, path, schema.ElastoDynBlade)
}

func (s *ElastoDynBlade) Format(path string) error {
	return format(s, path, schema.ElastoDynBlade)
}

type TableBlInpSt struct {
	Rows []TableBlInpStRow
}

type TableBlInpStRow struct {
	BlFract   float64 // -
	PitchAxis float64 // -
	StrcTwst  float64 // deg
	BMassDen  float64 // kg/m
	FlpStff   float64 // Nm^2
	EdgStff   float64 // Nm^2
}

func (t *TableBlInpSt) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableBlInpStRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableBlInpSt) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// ElastoDynTower
//------------------------------------------------------------------------------

type ElastoDynTower struct {
	Header1   Header       // ElastoDynTower Input File
	Title     Title        //
	Header2   Header       // Tower Parameters
	NTwInpSt  Int          // Number of input stations to specify tower geometry
	TwrFADmp1 Float        // Tower 1st fore-aft mode structural damping ratio
	TwrFADmp2 Float        // Tower 2nd fore-aft mode structural damping ratio
	TwrSSDmp1 Float        // Tower 1st side-to-side mode structural damping ratio
	TwrSSDmp2 Float        // Tower 2nd side-to-side mode structural damping ratio
	Header3   Header       // Tower Adjustment Factors
	FAStTunr1 Float        // Tower fore-aft modal stiffness tuner, 1st mode
	FAStTunr2 Float        // Tower fore-aft modal stiffness tuner, 2nd mode
	SSStTunr1 Float        // Tower side-to-side stiffness tuner, 1st mode
	SSStTunr2 Float        // Tower side-to-side stiffness tuner, 2nd mode
	AdjTwMa   Float        // Factor to adjust tower mass density
	AdjFASt   Float        // Factor to adjust tower fore-aft stiffness
	AdjSSSt   Float        // Factor to adjust tower side-to-side stiffness
	Header4   Header       // Distributed Tower Properties
	TwInpSt   TableTwInpSt //
	Header5   Header       // Tower Fore-Aft Mode Shapes
	TwFAM1Sh2 Float        // Mode 1, coefficient of x^2 term
	TwFAM1Sh3 Float        //       , coefficient of x^3 term
	TwFAM1Sh4 Float        //       , coefficient of x^4 term
	TwFAM1Sh5 Float        //       , coefficient of x^5 term
	TwFAM1Sh6 Float        //       , coefficient of x^6 term
	TwFAM2Sh2 Float        // Mode 2, coefficient of x^2 term
	TwFAM2Sh3 Float        //       , coefficient of x^3 term
	TwFAM2Sh4 Float        //       , coefficient of x^4 term
	TwFAM2Sh5 Float        //       , coefficient of x^5 term
	TwFAM2Sh6 Float        //       , coefficient of x^6 term
	Header6   Header       // Tower Side-To-Side Mode Shapes
	TwSSM1Sh2 Float        // Mode 1, coefficient of x^2 term
	TwSSM1Sh3 Float        //       , coefficient of x^3 term
	TwSSM1Sh4 Float        //       , coefficient of x^4 term
	TwSSM1Sh5 Float        //       , coefficient of x^5 term
	TwSSM1Sh6 Float        //       , coefficient of x^6 term
	TwSSM2Sh2 Float        // Mode 2, coefficient of x^2 term
	TwSSM2Sh3 Float        //       , coefficient of x^3 term
	TwSSM2Sh4 Float        //       , coefficient of x^4 term
	TwSSM2Sh5 Float        //       , coefficient of x^5 term
	TwSSM2Sh6 Float        //       , coefficient of x^6 term
}

func (s *ElastoDynTower) Parse(path string) error {
	return parse(s, path, schema.ElastoDynTower)
}

func (s *ElastoDynTower) Format(path string) error {
	return format(s, path, schema.ElastoDynTower)
}

type TableTwInpSt struct {
	Rows []TableTwInpStRow
}

type TableTwInpStRow struct {
	HtFract  float64 // -
	TMassDen float64 // kg/m
	TwFAStif float64 // Nm^2
	TwSSStif float64 // Nm^2
}

func (t *TableTwInpSt) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableTwInpStRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableTwInpSt) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// FreeVortexWake
//------------------------------------------------------------------------------

type FreeVortexWake struct {
	Text Text //
}

func (s *FreeVortexWake) Parse(path string) error {
	return parse(s, path, schema.FreeVortexWake)
}

func (s *FreeVortexWake) Format(path string) error {
	return format(s, path, schema.FreeVortexWake)
}

//------------------------------------------------------------------------------
// InflowWind
//------------------------------------------------------------------------------

type InflowWind struct {
	Header1        Header  // InflowWind Input File
	Title          Title   //
	Header2        Header  // General
	Echo           Bool    // Echo input data to <RootName>.ech
	WindType       Int     // switch for wind file type (1=steady; 2=uniform; 3=binary TurbSim FF; 4=binary Bladed-style FF; 5=HAWC format; 6=User defined; 7=native Bladed FF)
	PropagationDir Float   // Direction of wind propagation (meteorological rotation from aligned with X (positive rotates towards -Y) -- degrees) (not used for native Bladed format WindType=7)
	VFlowAng       Float   // Upflow angle (degrees) (not used for native Bladed format WindType=7)
	NWindVel       Int     // Number of points to output the wind velocity    (0 to 9)
	WindVxiList    Floats  // List of coordinates in the inertial X direction (m)
	WindVyiList    Floats  // List of coordinates in the inertial Y direction (m)
	WindVziList    Floats  // List of coordinates in the inertial Z direction (m)
	Header3        Header  // Parameters for Steady Wind Conditions [used only for WindType = 1]
	HWindSpeed     Float   // Horizontal wind speed
	RefHt          Float   // Reference height for horizontal wind speed
	PLExp          Float   // Power law exponent
	Header4        Header  // Parameters for Uniform wind file   [used only for WindType = 2]
	Filename_Uni   File    // Filename of time series data for uniform wind field
	RefHt_Uni      Float   // Reference height for horizontal wind speed
	RefLength      Float   // Reference length for linear horizontal and vertical sheer
	Header5        Header  // Parameters for Binary TurbSim Full-Field files   [used only for WindType = 3]
	FileName_BTS   Path    // Name of the Full field wind file to use (.bts)
	Header6        Header  // Parameters for Binary Bladed-style Full-Field files   [used only for WindType = 4 or WindType = 7]
	FileNameRoot   Path    // WindType=4: Rootname of the full-field wind file to use (.wnd, .sum); WindType=7: name of the intermediate file with wind scaling values
	TowerFile      Bool    // Have tower file (.twr) ignored when WindType = 7
	Header7        Header  // Parameters for HAWC-format binary files  [Only used with WindType = 5]
	FileName_u     Path    // name of the file containing the u-component fluctuating wind (.bin)
	FileName_v     Path    // name of the file containing the v-component fluctuating wind (.bin)
	FileName_w     Path    // name of the file containing the w-component fluctuating wind (.bin)
	Nx             Int     // number of grids in the x direction (in the 3 files above)
	Ny             Int     // number of grids in the y direction (in the 3 files above)
	Nz             Int     // number of grids in the z direction (in the 3 files above)
	Dx             Int     // distance (in meters) between points in the x direction
	Dy             Int     // distance (in meters) between points in the y direction
	Dz             Int     // distance (in meters) between points in the z direction
	RefHt_Hawc     Float   // reference height; the height (in meters) of the vertical center of the grid
	Header8        Header  // Scaling parameters for turbulence
	ScaleMethod    Int     // Turbulence scaling method [0 = none, 1 = direct scaling, 2 = calculate scaling factor based on a desired standard deviation]
	SFx            Float   // Turbulence scaling factor for the x direction [ScaleMethod=1]
	SFy            Float   // Turbulence scaling factor for the y direction [ScaleMethod=1]
	SFz            Float   // Turbulence scaling factor for the z direction [ScaleMethod=1]
	SigmaFx        Float   // Turbulence standard deviation to calculate scaling from in x direction [ScaleMethod=2]
	SigmaFy        Float   // Turbulence standard deviation to calculate scaling from in y direction [ScaleMethod=2]
	SigmaFz        Float   // Turbulence standard deviation to calculate scaling from in z direction [ScaleMethod=2]
	Header9        Header  // Mean wind profile parameters (added to HAWC-format files)
	URef           Float   // Mean u-component wind speed at the reference height
	WindProfile    Int     // Wind profile type (0=constant;1=logarithmic,2=power law)
	PLExp_Hawc     Float   // Power law exponent (used for PL wind profile type only)
	Z0             Float   // Surface roughness length (m) (used for LG wind profile type only)
	XOffset        Float   // Initial offset in +x direction (shift of wind box)
	Header10       Header  // OUTPUT
	SumPrint       Bool    // Print summary data to <RootName>.IfW.sum
	OutList        OutList // The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels
}

func (s *InflowWind) Parse(path string) error {
	return parse(s, path, schema.InflowWind)
}

func (s *InflowWind) Format(path string) error {
	return format(s, path, schema.InflowWind)
}

//------------------------------------------------------------------------------
// Main
//------------------------------------------------------------------------------

type Main struct {
	Header1     Header       // OpenFAST Input File
	Title       Title        //
	Header2     Header       // Simulation Control
	Echo        Bool         // Echo input data to <RootName>.ech
	AbortLevel  String       // Error level when simulation should abort (string)
	TMax        Float        // Total run time
	DT          Float        // Recommended module time step
	InterpOrder Int          // Interpolation order for input/output time history {1=linear, 2=quadratic}
	NumCrctn    Int          // Number of correction iterations {0=explicit calculation, i.e., no corrections}
	DT_UJac     Float        // Time between calls to get Jacobians
	UJacSclFact Float        // Scaling factor used in Jacobians
	Header3     Header       // Feature switches and flags
	CompElast   Int          // Compute structural dynamics
	CompInflow  Int          // Compute inflow wind velocities
	CompAero    Int          // Compute aerodynamic loads
	CompServo   Int          // Compute control and electrical-drive dynamics
	CompHydro   Int          // Compute hydrodynamic loads
	CompSub     Int          // Compute sub-structural dynamics
	CompMooring Int          // Compute mooring system
	CompIce     Int          // Compute ice loads
	MHK         Int          // MHK turbine type
	Header4     Header       // Environmental Conditions
	Gravity     Float        // Gravitational acceleration
	AirDens     Float        // Air density
	WtrDens     Float        // Water density
	KinVisc     Float        // Kinematic viscosity of working fluid
	SpdSound    Float        // Speed of sound in working fluid
	Patm        Float        // Atmospheric pressure [used only for an MHK turbine cavitation check]
	Pvap        Float        // Vapour pressure of working fluid [used only for an MHK turbine cavitation check]
	WtrDpth     Float        // Water depth
	MSL2SWL     Float        // Offset between still-water level and mean sea level [positive upward]
	Header5     Header       // Input Files
	EDFile      File         // Name of file containing ElastoDyn input parameters
	BDBldFile1  File         // Name of file containing BeamDyn input parameters for blade 1
	BDBldFile2  File         // Name of file containing BeamDyn input parameters for blade 2
	BDBldFile3  File         // Name of file containing BeamDyn input parameters for blade 3
	InflowFile  File         // Name of file containing inflow wind input parameters
	AeroFile    File         // Name of file containing aerodynamic input parameters
	ServoFile   File         // Name of file containing control and electrical-drive input parameters
	HydroFile   Path         // Name of file containing hydrodynamic input parameters
	SubFile     File         // Name of file containing sub-structural input parameters
	MooringFile Path         // Name of file containing mooring system input parameters
	IceFile     Path         // Name of file containing ice input parameters
	Header6     Header       // Output
	SumPrint    Bool         // Print summary data to '<RootName>.sum'
	SttsTime    Float        // Amount of time between screen status messages
	ChkptTime   Float        // Amount of time between creating checkpoint files for potential restart
	DT_Out      FloatDefault // Time step for tabular output (or "default")
	TStart      Float        // Time to begin tabular output
	OutFileFmt  Int          // Format for tabular (time-marching) output file
	TabDelim    Bool         // Use tab delimiters in text tabular output file?
	OutFmt      String       // Format used for text tabular output, excluding the time channel.  Resulting field should be 10 characters
	Header7     Header       // Linearization
	Linearize   Bool         // Linearization analysis
	CalcSteady  Bool         // Calculate a steady-state periodic operating point before linearization?
	TrimCase    Int          // Controller parameter to be trimmed
	TrimTol     Float        // Tolerance for the rotational speed convergence
	TrimGain    Float        // Proportional gain for the rotational speed error (>0) (rad/(rad/s) for yaw or pitch; Nm/(rad/s) for torque)
	Twr_Kdmp    Float        // Damping factor for the tower
	Bld_Kdmp    Float        // Damping factor for the blades
	NLinTimes   Int          // Number of times to linearize [>=1]
	LinTimes    Floats       // List of times at which to linearize [1 to NLinTimes] [used only when Linearize=True and CalcSteady=False]
	LinInputs   Int          // Inputs included in linearization
	LinOutputs  Int          // Outputs included in linearization
	LinOutJac   Bool         // Include full Jacobians in linearization output (for debug)
	LinOutMod   Bool         // Write module-level linearization output files in addition to output for full system?
	Header8     Header       // Visualization
	WrVTK       Int          // VTK visualization data output
	VTK_type    Int          // Type of VTK visualization data [unused if WrVTK=0]
	VTK_fields  Bool         // Write mesh fields to VTK data files? {true/false} [unused if WrVTK=0]
	VTK_fps     Float        // Frame rate for VTK output (frames per second) {will use closest integer multiple of DT} [used only if WrVTK=2 or WrVTK=3]
}

func (s *Main) Parse(path string) error {
	return parse(s, path, schema.Main)
}

func (s *Main) Format(path string) error {
	return format(s, path, schema.Main)
}

//------------------------------------------------------------------------------
// ServoDyn
//------------------------------------------------------------------------------

type ServoDyn struct {
	Header1      Header         // ServoDyn Input File
	Title        Title          //
	Header2      Header         // Simulation Control
	Echo         Bool           // Echo input data to "<rootname>.SD.ech"
	DT           FloatDefault   // Communication interval for controllers (s) (or "default")
	Header3      Header         // Pitch Control
	PCMode       Int            // Pitch control mode
	TPCOn        Float          // Time to enable active pitch control [unused when PCMode=0]
	TPitManS1    Float          // Time to start override pitch maneuver for blade 1 and end standard pitch control
	TPitManS2    Float          // Time to start override pitch maneuver for blade 2 and end standard pitch control
	TPitManS3    Float          // Time to start override pitch maneuver for blade 3 and end standard pitch control [unused for 2 blades]
	PitManRat1   Float          // Pitch rate at which override pitch maneuver heads toward final pitch angle for blade 1
	PitManRat2   Float          // Pitch rate at which override pitch maneuver heads toward final pitch angle for blade 2
	PitManRat3   Float          // Pitch rate at which override pitch maneuver heads toward final pitch angle for blade 3 [unused for 2 blades]
	BlPitchF1    Float          // Blade 1 final pitch for pitch maneuvers
	BlPitchF2    Float          // Blade 2 final pitch for pitch maneuvers
	BlPitchF3    Float          // Blade 3 final pitch for pitch maneuvers [unused for 2 blades]
	Header4      Header         // Generator And Torque Control
	VSContrl     Int            // Variable-speed control mode
	GenModel     Int            // Generator model [used only when VSContrl=0]
	GenEff       Float          // Generator efficiency [ignored by the Thevenin and user-defined generator models]
	GenTiStr     Bool           // Method to start the generator {T: timed using TimGenOn, F: generator speed using SpdGenOn}
	GenTiStp     Bool           // Method to stop the generator {T: timed using TimGenOf, F: when generator power = 0}
	SpdGenOn     Float          // Generator speed to turn on the generator for a startup (HSS speed) [used only when GenTiStr=False]
	TimGenOn     Float          // Time to turn on the generator for a startup (s) [used only when GenTiStr=True]
	TimGenOf     Float          // Time to turn off the generator (s) [used only when GenTiStp=True]
	Header5      Header         // Simple Variable-Speed Torque Control
	VS_RtGnSp    Float          // Rated generator speed for simple variable-speed generator control (HSS side) [used only when VSContrl=1]
	VS_RtTq      Float          // Rated generator torque/constant generator torque in Region 3 for simple variable-speed generator control (HSS side) [used only when VSContrl=1]
	VS_Rgn2K     Float          // Generator torque constant in Region 2 for simple variable-speed generator control (HSS side)  [used only when VSContrl=1]
	VS_SlPc      Float          // Rated generator slip percentage in Region 2 1/2 for simple variable-speed generator control [used only when VSContrl=1]
	Header6      Header         // Simple Induction Generator
	SIG_SlPc     Float          // Rated generator slip percentage [used only when VSContrl=0 and GenModel=1]
	SIG_SySp     Float          // Synchronous (zero-torque) generator speed [used only when VSContrl=0 and GenModel=1]
	SIG_RtTq     Float          // Rated torque [used only when VSContrl=0 and GenModel=1]
	SIG_PORt     Float          // Pull-out ratio [used only when VSContrl=0 and GenModel=1]
	Header7      Header         // Thevenin-Equivalent Induction Generator
	TEC_Freq     Float          // Line frequency [50 or 60] [used only when VSContrl=0 and GenModel=2]
	TEC_NPol     Float          // Number of poles [even integer > 0] [used only when VSContrl=0 and GenModel=2]
	TEC_SRes     Float          // Stator resistance [used only when VSContrl=0 and GenModel=2]
	TEC_RRes     Float          // Rotor resistance [used only when VSContrl=0 and GenModel=2]
	TEC_VLL      Float          // Line-to-line RMS voltage (volts) [used only when VSContrl=0 and GenModel=2]
	TEC_SLR      Float          // Stator leakage reactance [used only when VSContrl=0 and GenModel=2]
	TEC_RLR      Float          // Rotor leakage reactance [used only when VSContrl=0 and GenModel=2]
	TEC_MR       Float          // Magnetizing reactance [used only when VSContrl=0 and GenModel=2]
	Header8      Header         // High-Speed Shaft Brake
	HSSBrMode    Int            // HSS brake model
	THSSBrDp     Float          // Time to initiate deployment of the HSS brake
	HSSBrDT      Float          // Time for HSS-brake to reach full deployment once initiated (sec) [used only when HSSBrMode=1]
	HSSBrTqF     Float          // Fully deployed HSS-brake torque
	Header9      Header         // Nacelle-Yaw Control
	YCMode       Int            // Yaw control mode {0: none, 3: user-defined from routine UserYawCont, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL} (switch)
	TYCOn        Float          // Time to enable active yaw control [unused when YCMode=0]
	YawNeut      Float          // Neutral yaw position--yaw spring force is zero at this yaw
	YawSpr       Float          // Nacelle-yaw spring constant
	YawDamp      Float          // Nacelle-yaw damping constant
	TYawManS     Float          // Time to start override yaw maneuver and end standard yaw control
	YawManRat    Float          // Yaw maneuver rate (in absolute value)
	NacYawF      Float          // Final yaw angle for override yaw maneuvers
	Header10     Header         // Aerodynamic Flow Control
	AfCmode      Int            // Airfoil control mode
	AfC_Mean     Float          // Mean level for cosine cycling or steady value [used only with AfCmode==1]
	AfC_Amp      Float          // Amplitude for for cosine cycling of flap signal [used only with AfCmode==1]
	AfC_Phase    Float          // Phase relative to the blade azimuth (0 is vertical) for for cosine cycling of flap signal [used only with AfCmode==1]
	Header11     Header         // Structural Control
	NumBStC      Int            // Number of blade structural controllers (integer)
	BStCfiles    Paths          // Name of the files for blade structural controllers [unused when NumBStC==0]
	NumNStC      Int            // Number of nacelle structural controllers (integer)
	NStCfiles    Paths          // Name of the files for nacelle structural controllers [unused when NumNStC==0]
	NumTStC      Int            // Number of tower structural controllers (integer)
	TStCfiles    Paths          // Name of the files for tower structural controllers [unused when NumTStC==0]
	NumSStC      Int            // Number of substructure structural controllers (integer)
	SStCfiles    Paths          // Name of the files for substructure structural controllers [unused when NumSStC==0]
	Header12     Header         // Cable Control
	CCmode       Int            // Cable control mode {0: none, 4: user-defined from Simulink/Labview, 5: user-defined from Bladed-style DLL}
	Header13     Header         // BLADED INTERFACE [used only with Bladed Interface]
	DLL_FileName Path           // Name/location of the dynamic library {.dll [Windows] or .so [Linux]} in the Bladed-DLL format [used only with Bladed Interface]
	DLL_InFile   Path           // Name of input file sent to the DLL [used only with Bladed Interface]
	DLL_ProcName String         // Name of procedure in DLL to be called [case sensitive; used only with DLL Interface]
	DLL_DT       FloatDefault   // Communication interval for dynamic library (or "default") [used only with Bladed Interface]
	DLL_Ramp     Bool           // Whether a linear ramp should be used between DLL_DT time steps [introduces time shift when true] (flag) [used only with Bladed Interface]
	BPCutoff     Float          // Cutoff frequency for low-pass filter on blade pitch from DLL (Hz) [used only with Bladed Interface]
	NacYaw_North Float          // Reference yaw angle of the nacelle when the upwind end points due North [used only with Bladed Interface]
	Ptch_Cntrl   Float          // Record 28: Use individual pitch control {0: collective pitch; 1: individual pitch control} (switch) [used only with Bladed Interface]
	Ptch_SetPnt  Float          // Record  5: Below-rated pitch angle set-point [used only with Bladed Interface]
	Ptch_Min     Float          // Record  6: Minimum pitch angle [used only with Bladed Interface]
	Ptch_Max     Float          // Record  7: Maximum pitch angle [used only with Bladed Interface]
	PtchRate_Min Float          // Record  8: Minimum pitch rate (most negative value allowed) [used only with Bladed Interface]
	PtchRate_Max Float          // Record  9: Maximum pitch rate  [used only with Bladed Interface]
	Gain_OM      Float          // Record 16: Optimal mode gain [used only with Bladed Interface]
	GenSpd_MinOM Float          // Record 17: Minimum generator speed [used only with Bladed Interface]
	GenSpd_MaxOM Float          // Record 18: Optimal mode maximum speed [used only with Bladed Interface]
	GenSpd_Dem   Float          // Record 19: Demanded generator speed above rated [used only with Bladed Interface]
	GenTrq_Dem   Float          // Record 22: Demanded generator torque above rated [used only with Bladed Interface]
	GenPwr_Dem   Float          // Record 13: Demanded power [used only with Bladed Interface]
	Header14     Header         // Bladed Interface Torque-Speed Look-Up Table
	DLL_NumTrq   Int            // Record 26: No. of points in torque-speed look-up table {0 = none and use the optimal mode parameters; nonzero = ignore the optimal mode PARAMETERs by setting Record 16 to 0.0} [used only with Bladed Interface]
	GenSpdTrq    TableGenSpdTrq //
	Header15     Header         // Output
	SumPrint     Bool           // Print summary data to <RootName>.sum (flag) (currently unused)
	OutFile      Int            // Switch to determine where output will be placed: {1: in module output file only; 2: in glue code output file only; 3: both} (currently unused)
	TabDelim     Bool           // Use tab delimiters in text tabular output file? (flag) (currently unused)
	OutFmt       String         // Format used for text tabular output (except time).  Resulting field should be 10 characters. (quoted string) (currently unused)
	TStart       Float          // Time to begin tabular output (s) (currently unused)
	OutList      OutList        // The next line(s) contains a list of output parameters.  See OutListParameters.xlsx for a listing of available output channels, (-)
}

func (s *ServoDyn) Parse(path string) error {
	return parse(s, path, schema.ServoDyn)
}

func (s *ServoDyn) Format(path string) error {
	return format(s, path, schema.ServoDyn)
}

type TableGenSpdTrq struct {
	Rows []TableGenSpdTrqRow
}

type TableGenSpdTrqRow struct {
	GenSpd_TLU float64 // rpm
	GenTrq_TLU float64 // Nm
}

func (t *TableGenSpdTrq) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableGenSpdTrqRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableGenSpdTrq) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// SubDyn
//------------------------------------------------------------------------------

type SubDyn struct {
	Header1             Header               // SubDyn Input File
	Title               Title                //
	Header2             Header               // Simulation Control
	Echo                Bool                 // Echo input data to "<rootname>.SD.ech"
	SDdeltaT            FloatDefault         // Local Integration Step. If "default", the glue-code integration step will be used.
	IntMethod           Int                  // Integration Method [1/2/3/4 = RK4/AB4/ABM4/AM2].
	SttcSolve           Bool                 // Solve dynamics about static equilibrium point
	GuyanLoadCorrection Bool                 // Include extra moment from lever arm at interface and rotate FEM for floating.
	Header3             Header               // Fea And Craig-Bampton Parameters
	FEMMod              Int                  // FEM switch: element model in the FEM. [1= Euler-Bernoulli(E-B);  2=Tapered E-B (unavailable);  3= 2-node Timoshenko;  4= 2-node tapered Timoshenko (unavailable)]
	NDiv                Int                  // Number of sub-elements per member
	CBMod               Bool                 // [T/F] If True perform C-B reduction, else full FEM dofs will be retained. If True, select Nmodes to retain in C-B reduced system.
	Nmodes              Int                  // Number of internal modes to retain (ignored if CBMod=False). If Nmodes=0 --> Guyan Reduction.
	JDampings           Floats               // Damping Ratios for each retained mode (% of critical) If Nmodes>0, list Nmodes structural damping ratios for each retained mode (% of critical), or a single damping ratio to be applied to all retained modes. (last entered value will be used for all remaining modes).
	GuyanDampMod        Int                  // Guyan damping {0=none, 1=Rayleigh Damping, 2=user specified 6x6 matrix}
	RayleighDamp        Floats               // Mass and stiffness proportional damping  coefficients (Rayleigh Damping) [only if GuyanDampMod=1]
	GuyanDampSize       Int                  // Guyan damping matrix (6x6) [only if GuyanDampMod=2]
	GuyanDampMatrix     TableGuyanDampMatrix //
	Header4             Header               // STRUCTURE JOINTS: joints connect structure members (~Hydrodyn Input File)
	NJoints             Int                  // Number of joints
	Joints              TableJoints          //
	Header5             Header               // Base Reaction Joints: 1/0 For Locked/Free Dof @ Each Reaction Node
	NReact              Int                  // Number of Joints with reaction forces; be sure to remove all rigid motion DOFs of the structure  (else det([K])=[0])
	React               TableReact           //
	Header6             Header               // INTERFACE JOINTS: 1/0 for Locked (to the TP)/Free DOF @each Interface Joint (only Locked-to-TP implemented thus far (=rigid TP))
	NInterf             Int                  // Number of interface joints locked to the Transition Piece (TP):  be sure to remove all rigid motion dofs
	Interf              TableInterf          //
	Header7             Header               // Members
	NMembers            Int                  // Number of frame members
	Members             TableMembers         //
	Header8             Header               // Member X-Section Property Data 1/2 [Isotropic Material For Now: Use This Table For Circular-Tubular Elements]
	NPropSets           Int                  // Number of structurally unique x-sections (i.e. how many groups of X-sectional properties are utilized throughout all of the members)
	PropSets            TablePropSets        //
	Header9             Header               // Member X-Section Property Data 2/2 [Isotropic Material For Now: Use This Table If Any Section Other Than Circular, However Provide Cosm(I,J) Below]
	NXPropSets          Int                  // Number of structurally unique non-circular x-sections (if 0 the following table is ignored)
	XPropSets           TableXPropSets       //
	Header10            Header               // Cable Properties
	NCablePropSets      Int                  // Number of cable cable properties
	CablePropSets       TableCablePropSets   //
	Header11            Header               // Rigid Link Properties
	NRigidPropSets      Int                  // Number of rigid link properties
	RigidPropSets       TableRigidPropSets   //
	Header12            Header               // Member Cosine Matrices Cosm(I,J)
	NCOSMs              Int                  // Number of unique cosine matrices (i.e., of unique member alignments including principal axis rotations); ignored if NXPropSets=0   or 9999 in any element below
	COSMs               TableCOSMs           //
	Header13            Header               // Joint Additional Concentrated Masses
	NCmass              Int                  // Number of joints with concentrated masses; Global Coordinate System
	Cmass               TableCmass           //
	Header14            Header               // Output: Summary & Outfile
	SumPrint            Bool                 // Output a Summary File .It contains: matrices K,M  and C-B reduced M_BB, M-BM, K_BB, K_MM(OMG^2), PHI_R, PHI_L. It can also contain COSMs if requested.
	OutCBModes          Int                  // Output Guyan and Craig-Bampton modes {0: No output, 1: JSON output}
	OutFEMModes         Int                  // Output first 30 FEM modes {0: No output, 1: JSON output}
	OutCOSM             Bool                 // Output cosine matrices with the selected output member forces
	OutAll              Bool                 // [T/F] Output all members' end forces
	OutSwtch            Int                  // [1/2/3] Output requested channels to: 1=<rootname>.SD.out;  2=<rootname>.out (generated by FAST);  3=both files.
	TabDelim            Bool                 // Generate a tab-delimited output in the <rootname>.SD.out file
	OutDec              Int                  // Decimation of output in the <rootname>.SD.out file
	OutFmt              String               // Output format for numerical results in the <rootname>.SD.out file
	OutSFmt             String               // Output format for header strings in the <rootname>.SD.out file
	Header15            Header               // Member Output List
	NMOutputs           Int                  // Number of members whose forces/displacements/velocities/accelerations will be output [Must be <= 9].
	MOutputs            TableMOutputs        //
	Header16            Header               // SSOutList: The next line(s) contains a list of output parameters that will be output in <rootname>.SD.out or <rootname>.out.
	OutList             OutList2             //
}

func (s *SubDyn) Parse(path string) error {
	return parse(s, path, schema.SubDyn)
}

func (s *SubDyn) Format(path string) error {
	return format(s, path, schema.SubDyn)
}

type TableGuyanDampMatrix struct {
	Rows []TableGuyanDampMatrixRow
}

type TableGuyanDampMatrixRow struct {
	C1 float64 //
	C2 float64 //
	C3 float64 //
	C4 float64 //
	C5 float64 //
	C6 float64 //
}

func (t *TableGuyanDampMatrix) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableGuyanDampMatrixRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableGuyanDampMatrix) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableJoints struct {
	Rows []TableJointsRow
}

type TableJointsRow struct {
	JointID    int     // -
	JointXss   float64 // m
	JointYss   float64 // m
	JointZss   float64 // m
	JointType  int     // -
	JointDirX  float64 // -
	JointDirY  float64 // -
	JointDirZ  float64 // -
	JointStiff float64 // Nm/rad
}

func (t *TableJoints) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableJointsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableJoints) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableReact struct {
	Rows []TableReactRow
}

type TableReactRow struct {
	RJointID int    // -
	RctTDXss int    // flag
	RctTDYss int    // flag
	RctTDZss int    // flag
	RctRDXss int    // flag
	RctRDYss int    // flag
	RctRDZss int    // flag
	SSIfile  string // string
}

func (t *TableReact) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableReactRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableReact) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableInterf struct {
	Rows []TableInterfRow
}

type TableInterfRow struct {
	IJointID int // -
	ItfTDXss int // flag
	ItfTDYss int // flag
	ItfTDZss int // flag
	ItfRDXss int // flag
	ItfRDYss int // flag
	ItfRDZss int // flag
}

func (t *TableInterf) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableInterfRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableInterf) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableMembers struct {
	Rows []TableMembersRow
}

type TableMembersRow struct {
	MemberID    int // -
	MJointID1   int // -
	MJointID2   int // -
	MPropSetID1 int // -
	MPropSetID2 int // -
	MType       int // -
	COSMID      int // -
}

func (t *TableMembers) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableMembersRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableMembers) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TablePropSets struct {
	Rows []TablePropSetsRow
}

type TablePropSetsRow struct {
	PropSetID float64 // -
	YoungE    float64 // N/m2
	ShearG    float64 // N/m2
	MatDens   float64 // kg/m3
	XsecD     float64 // m
	XsecT     float64 // m
}

func (t *TablePropSets) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TablePropSetsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TablePropSets) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableXPropSets struct {
	Rows []TableXPropSetsRow
}

type TableXPropSetsRow struct {
	PropSetID int     // -
	YoungE    float64 // N/m2
	ShearG    float64 // N/m2
	MatDens   float64 // kg/m3
	XsecA     float64 // m2
	XsecAsx   float64 // m2
	XsecAsy   float64 // m2
	XsecJxx   float64 // m4
	XsecJyy   float64 // m4
	XsecJ0    float64 // m4
}

func (t *TableXPropSets) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableXPropSetsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableXPropSets) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableCablePropSets struct {
	Rows []TableCablePropSetsRow
}

type TableCablePropSetsRow struct {
	PropSetID   int     // -
	EA          float64 // N
	MatDens     float64 // kg/m
	T0          float64 // N
	CtrlChannel int     // -
}

func (t *TableCablePropSets) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableCablePropSetsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableCablePropSets) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableRigidPropSets struct {
	Rows []TableRigidPropSetsRow
}

type TableRigidPropSetsRow struct {
	PropSetID int     // -
	EA        float64 // N
	MatDens   float64 // kg/m
}

func (t *TableRigidPropSets) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableRigidPropSetsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableRigidPropSets) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableCOSMs struct {
	Rows []TableCOSMsRow
}

type TableCOSMsRow struct {
	COSMID int     // -
	COSM11 float64 // -
	COSM12 float64 // -
	COSM13 float64 // -
	COSM21 float64 // -
	COSM22 float64 // -
	COSM23 float64 // -
	COSM31 float64 // -
	COSM32 float64 // -
	COSM33 float64 // -
}

func (t *TableCOSMs) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableCOSMsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableCOSMs) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableCmass struct {
	Rows []TableCmassRow
}

type TableCmassRow struct {
	CMJointID int     // -
	JMass     float64 // kg
	JMXX      float64 // kg*m^2
	JMYY      float64 // kg*m^2
	JMZZ      float64 // kg*m^2
	JMXY      float64 // kg*m^2
	JMXZ      float64 // kg*m^2
	JMYZ      float64 // kg*m^2
	MCGX      float64 // m
	MCGY      float64 // m
	MCGZ      float64 // m
}

func (t *TableCmass) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableCmassRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableCmass) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

type TableMOutputs struct {
	Rows []TableMOutputsRow
}

type TableMOutputsRow struct {
	MemberID int // -
	NOutCnt  int // -
	NodeCnt  int // -
}

func (t *TableMOutputs) Parse(s *schema.Schema, lines []string, num int) ([]string, error) {
	if n := num + s.TableHeaderSize; len(lines) < n {
		return nil, fmt.Errorf("insufficient lines, need %d", n)
	}
	lines = lines[s.TableHeaderSize:]
	t.Rows = make([]TableMOutputsRow, num)
	for i, line := range lines[:num] {
		if err := lineToStruct(&t.Rows[i], line); err != nil {
			return nil, err
		}
	}
	return lines[num:], nil
}

func (t TableMOutputs) Format(s *schema.Schema, w *bytes.Buffer) {
	if s.TableHeaderSize > 0 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", c.Name)
		}
		w.WriteString("\n")
	}
	if s.TableHeaderSize > 1 {
		for _, c := range s.TableColumns {
			fmt.Fprintf(w, " %14s", "("+c.Unit+")")
		}
		w.WriteString("\n")
	}
	for _, row := range t.Rows {
		structToLine(row, w)
	}
}

//------------------------------------------------------------------------------
// TailFin
//------------------------------------------------------------------------------

type TailFin struct {
	Text Text //
}

func (s *TailFin) Parse(path string) error {
	return parse(s, path, schema.TailFin)
}

func (s *TailFin) Format(path string) error {
	return format(s, path, schema.TailFin)
}

//------------------------------------------------------------------------------
// TextFile
//------------------------------------------------------------------------------

type TextFile struct {
	Text Text //
}

func (s *TextFile) Parse(path string) error {
	return parse(s, path, schema.TextFile)
}

func (s *TextFile) Format(path string) error {
	return format(s, path, schema.TextFile)
}

//------------------------------------------------------------------------------
// UniformWind
//------------------------------------------------------------------------------

type UniformWind struct {
	Text Text //
}

func (s *UniformWind) Parse(path string) error {
	return parse(s, path, schema.UniformWind)
}

func (s *UniformWind) Format(path string) error {
	return format(s, path, schema.UniformWind)
}

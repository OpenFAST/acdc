export namespace main {
	
	export class Paths {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: string[];
	    FileType: string;
	    Condensed: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Paths(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Desc = source["Desc"];
	        this.Line = source["Line"];
	        this.Value = source["Value"];
	        this.FileType = source["FileType"];
	        this.Condensed = source["Condensed"];
	    }
	}
	export class Integer {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: number;
	    Size: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Integer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Desc = source["Desc"];
	        this.Line = source["Line"];
	        this.Value = source["Value"];
	        this.Size = source["Size"];
	    }
	}
	export class Path {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: string;
	    FileType: string;
	    Root: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Path(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Desc = source["Desc"];
	        this.Line = source["Line"];
	        this.Value = source["Value"];
	        this.FileType = source["FileType"];
	        this.Root = source["Root"];
	    }
	}
	export class AeroDyn {
	    Name: string;
	    Type: string;
	    Text: string;
	    OLAFInputFileName: Path;
	    NumAFfiles: Integer;
	    AFNames: Paths;
	    ADBlFile1: Path;
	    ADBlFile2: Path;
	    ADBlFile3: Path;
	    TFinFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new AeroDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.OLAFInputFileName = this.convertValues(source["OLAFInputFileName"], Path);
	        this.NumAFfiles = this.convertValues(source["NumAFfiles"], Integer);
	        this.AFNames = this.convertValues(source["AFNames"], Paths);
	        this.ADBlFile1 = this.convertValues(source["ADBlFile1"], Path);
	        this.ADBlFile2 = this.convertValues(source["ADBlFile2"], Path);
	        this.ADBlFile3 = this.convertValues(source["ADBlFile3"], Path);
	        this.TFinFile = this.convertValues(source["TFinFile"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AeroDyn14 {
	    Name: string;
	    Type: string;
	    Text: string;
	    NumFoil: Integer;
	    FoilNm: Paths;
	
	    static createFrom(source: any = {}) {
	        return new AeroDyn14(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.NumFoil = this.convertValues(source["NumFoil"], Integer);
	        this.FoilNm = this.convertValues(source["FoilNm"], Paths);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AirfoilInfo {
	    Name: string;
	    Type: string;
	    Text: string;
	    BL_File: Path;
	
	    static createFrom(source: any = {}) {
	        return new AirfoilInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.BL_File = this.convertValues(source["BL_File"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Condition {
	    ID: number;
	    WindSpeed: number;
	    RotorSpeed: number;
	    BladePitch: number;
	
	    static createFrom(source: any = {}) {
	        return new Condition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.WindSpeed = source["WindSpeed"];
	        this.RotorSpeed = source["RotorSpeed"];
	        this.BladePitch = source["BladePitch"];
	    }
	}
	export class Range {
	    Min: number;
	    Max: number;
	    Num: number;
	
	    static createFrom(source: any = {}) {
	        return new Range(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Min = source["Min"];
	        this.Max = source["Max"];
	        this.Num = source["Num"];
	    }
	}
	export class Case {
	    ID: number;
	    Name: string;
	    IncludeAero: boolean;
	    RotorSpeedRange: Range;
	    WindSpeedRange: Range;
	    Curve: Condition[];
	    OperatingPoints: Condition[];
	
	    static createFrom(source: any = {}) {
	        return new Case(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.IncludeAero = source["IncludeAero"];
	        this.RotorSpeedRange = this.convertValues(source["RotorSpeedRange"], Range);
	        this.WindSpeedRange = this.convertValues(source["WindSpeedRange"], Range);
	        this.Curve = this.convertValues(source["Curve"], Condition);
	        this.OperatingPoints = this.convertValues(source["OperatingPoints"], Condition);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Analysis {
	    Cases: Case[];
	
	    static createFrom(source: any = {}) {
	        return new Analysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Cases = this.convertValues(source["Cases"], Case);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BeamDyn {
	    Name: string;
	    Type: string;
	    Text: string;
	    BldFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new BeamDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.BldFile = this.convertValues(source["BldFile"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Bool {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Bool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Desc = source["Desc"];
	        this.Line = source["Line"];
	        this.Value = source["Value"];
	    }
	}
	
	
	export class Config {
	    RecentProjects: string[];
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RecentProjects = source["RecentProjects"];
	    }
	}
	export class Real {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: number;
	
	    static createFrom(source: any = {}) {
	        return new Real(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Desc = source["Desc"];
	        this.Line = source["Line"];
	        this.Value = source["Value"];
	    }
	}
	export class ElastoDyn {
	    Name: string;
	    Type: string;
	    Text: string;
	    FlapDOF1: Bool;
	    FlapDOF2: Bool;
	    EdgeDOF: Bool;
	    TeetDOF: Bool;
	    DrTrDOF: Bool;
	    GenDOF: Bool;
	    YawDOF: Bool;
	    TwFADOF1: Bool;
	    TwFADOF2: Bool;
	    TwSSDOF1: Bool;
	    TwSSDOF2: Bool;
	    BlPitch1: Real;
	    BlPitch2: Real;
	    BlPitch3: Real;
	    RotSpeed: Real;
	    NumBl: Integer;
	    BldFile1: Path;
	    BldFile2: Path;
	    BldFile3: Path;
	    TwrFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new ElastoDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.FlapDOF1 = this.convertValues(source["FlapDOF1"], Bool);
	        this.FlapDOF2 = this.convertValues(source["FlapDOF2"], Bool);
	        this.EdgeDOF = this.convertValues(source["EdgeDOF"], Bool);
	        this.TeetDOF = this.convertValues(source["TeetDOF"], Bool);
	        this.DrTrDOF = this.convertValues(source["DrTrDOF"], Bool);
	        this.GenDOF = this.convertValues(source["GenDOF"], Bool);
	        this.YawDOF = this.convertValues(source["YawDOF"], Bool);
	        this.TwFADOF1 = this.convertValues(source["TwFADOF1"], Bool);
	        this.TwFADOF2 = this.convertValues(source["TwFADOF2"], Bool);
	        this.TwSSDOF1 = this.convertValues(source["TwSSDOF1"], Bool);
	        this.TwSSDOF2 = this.convertValues(source["TwSSDOF2"], Bool);
	        this.BlPitch1 = this.convertValues(source["BlPitch1"], Real);
	        this.BlPitch2 = this.convertValues(source["BlPitch2"], Real);
	        this.BlPitch3 = this.convertValues(source["BlPitch3"], Real);
	        this.RotSpeed = this.convertValues(source["RotSpeed"], Real);
	        this.NumBl = this.convertValues(source["NumBl"], Integer);
	        this.BldFile1 = this.convertValues(source["BldFile1"], Path);
	        this.BldFile2 = this.convertValues(source["BldFile2"], Path);
	        this.BldFile3 = this.convertValues(source["BldFile3"], Path);
	        this.TwrFile = this.convertValues(source["TwrFile"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EvalStatus {
	    ID: number;
	    State: string;
	    Progress: number;
	    Error: string;
	
	    static createFrom(source: any = {}) {
	        return new EvalStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.State = source["State"];
	        this.Progress = source["Progress"];
	        this.Error = source["Error"];
	    }
	}
	export class Exec {
	    Path: string;
	    Version: string;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Exec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Version = source["Version"];
	        this.Valid = source["Valid"];
	    }
	}
	export class StControl {
	    Name: string;
	    Type: string;
	    Text: string;
	    PrescribedForcesFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new StControl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.PrescribedForcesFile = this.convertValues(source["PrescribedForcesFile"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Misc {
	    Name: string;
	    Type: string;
	    Text: string;
	
	    static createFrom(source: any = {}) {
	        return new Misc(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	    }
	}
	export class OLAF {
	    Name: string;
	    Type: string;
	    Text: string;
	    PrescribedCircFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new OLAF(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.PrescribedCircFile = this.convertValues(source["PrescribedCircFile"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InflowWind {
	    Name: string;
	    Type: string;
	    Text: string;
	    WindType: Integer;
	    PropagationDir: Real;
	    VFlowAng: Real;
	    HWindSpeed: Real;
	    PLExp: Real;
	
	    static createFrom(source: any = {}) {
	        return new InflowWind(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.WindType = this.convertValues(source["WindType"], Integer);
	        this.PropagationDir = this.convertValues(source["PropagationDir"], Real);
	        this.VFlowAng = this.convertValues(source["VFlowAng"], Real);
	        this.HWindSpeed = this.convertValues(source["HWindSpeed"], Real);
	        this.PLExp = this.convertValues(source["PLExp"], Real);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ServoDyn {
	    Name: string;
	    Type: string;
	    Text: string;
	    NumBStC: Integer;
	    BStCfiles: Paths;
	    NumNStC: Integer;
	    NStCfiles: Paths;
	    NumTStC: Integer;
	    TStCfiles: Paths;
	    NumSStC: Integer;
	    SStCfiles: Paths;
	
	    static createFrom(source: any = {}) {
	        return new ServoDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.NumBStC = this.convertValues(source["NumBStC"], Integer);
	        this.BStCfiles = this.convertValues(source["BStCfiles"], Paths);
	        this.NumNStC = this.convertValues(source["NumNStC"], Integer);
	        this.NStCfiles = this.convertValues(source["NStCfiles"], Paths);
	        this.NumTStC = this.convertValues(source["NumTStC"], Integer);
	        this.TStCfiles = this.convertValues(source["TStCfiles"], Paths);
	        this.NumSStC = this.convertValues(source["NumSStC"], Integer);
	        this.SStCfiles = this.convertValues(source["SStCfiles"], Paths);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HydroDyn {
	    Name: string;
	    Type: string;
	    Text: string;
	    PotFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new HydroDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.PotFile = this.convertValues(source["PotFile"], Path);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Reals {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: number[];
	
	    static createFrom(source: any = {}) {
	        return new Reals(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Desc = source["Desc"];
	        this.Line = source["Line"];
	        this.Value = source["Value"];
	    }
	}
	export class Main {
	    Name: string;
	    Type: string;
	    Text: string;
	    DT: Real;
	    CompElast: Integer;
	    CompInflow: Integer;
	    CompAero: Integer;
	    CompServo: Integer;
	    CompHydro: Integer;
	    CompSub: Integer;
	    CompMooring: Integer;
	    CompIce: Integer;
	    MHK: Integer;
	    EDFile: Path;
	    BDBldFile1: Path;
	    BDBldFile2: Path;
	    BDBldFile3: Path;
	    InflowFile: Path;
	    AeroFile: Path;
	    ServoFile: Path;
	    HydroFile: Path;
	    SubFile: Path;
	    MooringFile: Path;
	    IceFile: Path;
	    Linearize: Bool;
	    CalcSteady: Bool;
	    TrimCase: Integer;
	    TrimTol: Real;
	    TrimGain: Real;
	    Twr_Kdmp: Real;
	    Bld_Kdmp: Real;
	    NLinTimes: Integer;
	    LinTimes: Reals;
	    LinInputs: Integer;
	    LinOutputs: Integer;
	    LinOutJac: Bool;
	    LinOutMod: Bool;
	
	    static createFrom(source: any = {}) {
	        return new Main(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Text = source["Text"];
	        this.DT = this.convertValues(source["DT"], Real);
	        this.CompElast = this.convertValues(source["CompElast"], Integer);
	        this.CompInflow = this.convertValues(source["CompInflow"], Integer);
	        this.CompAero = this.convertValues(source["CompAero"], Integer);
	        this.CompServo = this.convertValues(source["CompServo"], Integer);
	        this.CompHydro = this.convertValues(source["CompHydro"], Integer);
	        this.CompSub = this.convertValues(source["CompSub"], Integer);
	        this.CompMooring = this.convertValues(source["CompMooring"], Integer);
	        this.CompIce = this.convertValues(source["CompIce"], Integer);
	        this.MHK = this.convertValues(source["MHK"], Integer);
	        this.EDFile = this.convertValues(source["EDFile"], Path);
	        this.BDBldFile1 = this.convertValues(source["BDBldFile1"], Path);
	        this.BDBldFile2 = this.convertValues(source["BDBldFile2"], Path);
	        this.BDBldFile3 = this.convertValues(source["BDBldFile3"], Path);
	        this.InflowFile = this.convertValues(source["InflowFile"], Path);
	        this.AeroFile = this.convertValues(source["AeroFile"], Path);
	        this.ServoFile = this.convertValues(source["ServoFile"], Path);
	        this.HydroFile = this.convertValues(source["HydroFile"], Path);
	        this.SubFile = this.convertValues(source["SubFile"], Path);
	        this.MooringFile = this.convertValues(source["MooringFile"], Path);
	        this.IceFile = this.convertValues(source["IceFile"], Path);
	        this.Linearize = this.convertValues(source["Linearize"], Bool);
	        this.CalcSteady = this.convertValues(source["CalcSteady"], Bool);
	        this.TrimCase = this.convertValues(source["TrimCase"], Integer);
	        this.TrimTol = this.convertValues(source["TrimTol"], Real);
	        this.TrimGain = this.convertValues(source["TrimGain"], Real);
	        this.Twr_Kdmp = this.convertValues(source["Twr_Kdmp"], Real);
	        this.Bld_Kdmp = this.convertValues(source["Bld_Kdmp"], Real);
	        this.NLinTimes = this.convertValues(source["NLinTimes"], Integer);
	        this.LinTimes = this.convertValues(source["LinTimes"], Reals);
	        this.LinInputs = this.convertValues(source["LinInputs"], Integer);
	        this.LinOutputs = this.convertValues(source["LinOutputs"], Integer);
	        this.LinOutJac = this.convertValues(source["LinOutJac"], Bool);
	        this.LinOutMod = this.convertValues(source["LinOutMod"], Bool);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Files {
	    Main: Main[];
	    ElastoDyn: ElastoDyn[];
	    BeamDyn: BeamDyn[];
	    AeroDyn: AeroDyn[];
	    AeroDyn14: AeroDyn14[];
	    HydroDyn: HydroDyn[];
	    ServoDyn: ServoDyn[];
	    InflowWind: InflowWind[];
	    OLAF: OLAF[];
	    Misc: Misc[];
	    StControl: StControl[];
	    AirfoilInfo: AirfoilInfo[];
	
	    static createFrom(source: any = {}) {
	        return new Files(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Main = this.convertValues(source["Main"], Main);
	        this.ElastoDyn = this.convertValues(source["ElastoDyn"], ElastoDyn);
	        this.BeamDyn = this.convertValues(source["BeamDyn"], BeamDyn);
	        this.AeroDyn = this.convertValues(source["AeroDyn"], AeroDyn);
	        this.AeroDyn14 = this.convertValues(source["AeroDyn14"], AeroDyn14);
	        this.HydroDyn = this.convertValues(source["HydroDyn"], HydroDyn);
	        this.ServoDyn = this.convertValues(source["ServoDyn"], ServoDyn);
	        this.InflowWind = this.convertValues(source["InflowWind"], InflowWind);
	        this.OLAF = this.convertValues(source["OLAF"], OLAF);
	        this.Misc = this.convertValues(source["Misc"], Misc);
	        this.StControl = this.convertValues(source["StControl"], StControl);
	        this.AirfoilInfo = this.convertValues(source["AirfoilInfo"], AirfoilInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Info {
	    Date: string;
	    Path: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Date = source["Date"];
	        this.Path = source["Path"];
	    }
	}
	
	
	export class Model {
	    HasAero: boolean;
	    ImportedPaths: string[];
	    Files?: Files;
	    Notes: string[];
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HasAero = source["HasAero"];
	        this.ImportedPaths = source["ImportedPaths"];
	        this.Files = this.convertValues(source["Files"], Files);
	        this.Notes = source["Notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	

}


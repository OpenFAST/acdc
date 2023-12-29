export namespace diagram {
	
	export class Point {
	    OpPtID: number;
	    ModeID: number;
	    RotSpeed: number;
	    WindSpeed: number;
	    NaturalFreqHz: number;
	    DampedFreqHz: number;
	    DampingRatio: number;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OpPtID = source["OpPtID"];
	        this.ModeID = source["ModeID"];
	        this.RotSpeed = source["RotSpeed"];
	        this.WindSpeed = source["WindSpeed"];
	        this.NaturalFreqHz = source["NaturalFreqHz"];
	        this.DampedFreqHz = source["DampedFreqHz"];
	        this.DampingRatio = source["DampingRatio"];
	    }
	}
	export class Line {
	    ID: number;
	    Label: string;
	    Points: Point[];
	    Hide: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Line(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Label = source["Label"];
	        this.Points = this.convertValues(source["Points"], Point);
	        this.Hide = source["Hide"];
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
	export class Diagram {
	    HasWind: boolean;
	    RotSpeeds: number[];
	    WindSpeeds: number[];
	    Lines: Line[];
	
	    static createFrom(source: any = {}) {
	        return new Diagram(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HasWind = source["HasWind"];
	        this.RotSpeeds = source["RotSpeeds"];
	        this.WindSpeeds = source["WindSpeeds"];
	        this.Lines = this.convertValues(source["Lines"], Line);
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

export namespace lin {
	
	export class Mode {
	    ID: number;
	    OP: number;
	    EigenValueReal: number;
	    EigenValueImag: number;
	    NaturalFreqRaw: number;
	    NaturalFreqHz: number;
	    DampedFreqRaw: number;
	    DampedFreqHz: number;
	    DampingRatio: number;
	    Magnitudes: number[];
	    Phases: number[];
	
	    static createFrom(source: any = {}) {
	        return new Mode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.OP = source["OP"];
	        this.EigenValueReal = source["EigenValueReal"];
	        this.EigenValueImag = source["EigenValueImag"];
	        this.NaturalFreqRaw = source["NaturalFreqRaw"];
	        this.NaturalFreqHz = source["NaturalFreqHz"];
	        this.DampedFreqRaw = source["DampedFreqRaw"];
	        this.DampedFreqHz = source["DampedFreqHz"];
	        this.DampingRatio = source["DampingRatio"];
	        this.Magnitudes = source["Magnitudes"];
	        this.Phases = source["Phases"];
	    }
	}
	export class EigenResults {
	    Modes: Mode[];
	    // Go type: mat
	    EigenVectors?: any;
	
	    static createFrom(source: any = {}) {
	        return new EigenResults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Modes = this.convertValues(source["Modes"], Mode);
	        this.EigenVectors = this.convertValues(source["EigenVectors"], null);
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
	export class OPOrder {
	    Num: number;
	    NumFixed: number;
	    NumRotating: number;
	    NumTriplets: number;
	    Indices: number[];
	    Triplets: number[][];
	
	    static createFrom(source: any = {}) {
	        return new OPOrder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Num = source["Num"];
	        this.NumFixed = source["NumFixed"];
	        this.NumRotating = source["NumRotating"];
	        this.NumTriplets = source["NumTriplets"];
	        this.Indices = source["Indices"];
	        this.Triplets = source["Triplets"];
	    }
	}
	export class MBC {
	    RotSpeed: number;
	    WindSpeed: number;
	    DescStates: string[];
	    Azimuths: number[];
	    OrderX: OPOrder;
	    OrderX2: OPOrder;
	    OrderX2dot: OPOrder;
	    OrderX1: OPOrder;
	    OrderU: OPOrder;
	    OrderY: OPOrder;
	    OrderEigen: OPOrder;
	    DOFsEigen: string[];
	
	    static createFrom(source: any = {}) {
	        return new MBC(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RotSpeed = source["RotSpeed"];
	        this.WindSpeed = source["WindSpeed"];
	        this.DescStates = source["DescStates"];
	        this.Azimuths = source["Azimuths"];
	        this.OrderX = this.convertValues(source["OrderX"], OPOrder);
	        this.OrderX2 = this.convertValues(source["OrderX2"], OPOrder);
	        this.OrderX2dot = this.convertValues(source["OrderX2dot"], OPOrder);
	        this.OrderX1 = this.convertValues(source["OrderX1"], OPOrder);
	        this.OrderU = this.convertValues(source["OrderU"], OPOrder);
	        this.OrderY = this.convertValues(source["OrderY"], OPOrder);
	        this.OrderEigen = this.convertValues(source["OrderEigen"], OPOrder);
	        this.DOFsEigen = source["DOFsEigen"];
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
	
	
	export class OPResults {
	    Name: string;
	    MBC?: MBC;
	    EigRes?: EigenResults;
	
	    static createFrom(source: any = {}) {
	        return new OPResults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.MBC = this.convertValues(source["MBC"], MBC);
	        this.EigRes = this.convertValues(source["EigRes"], EigenResults);
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
	export class AeroDyn {
	    Name: string;
	    Type: string;
	    Lines: string[];
	    WakeMod: Integer;
	    AFAeroMod: Integer;
	    TwrPotent: Integer;
	    TwrShadow: Integer;
	    FrozenWake: Bool;
	    SkewMod: Integer;
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
	        this.Lines = source["Lines"];
	        this.WakeMod = this.convertValues(source["WakeMod"], Integer);
	        this.AFAeroMod = this.convertValues(source["AFAeroMod"], Integer);
	        this.TwrPotent = this.convertValues(source["TwrPotent"], Integer);
	        this.TwrShadow = this.convertValues(source["TwrShadow"], Integer);
	        this.FrozenWake = this.convertValues(source["FrozenWake"], Bool);
	        this.SkewMod = this.convertValues(source["SkewMod"], Integer);
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
	    Lines: string[];
	    NumFoil: Integer;
	    FoilNm: Paths;
	
	    static createFrom(source: any = {}) {
	        return new AeroDyn14(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
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
	    Lines: string[];
	    BL_File: Path;
	
	    static createFrom(source: any = {}) {
	        return new AirfoilInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
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
	    CutIn: number;
	    Rated: number;
	    CutOut: number;
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
	        this.CutIn = source["CutIn"];
	        this.Rated = source["Rated"];
	        this.CutOut = source["CutOut"];
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
	    Lines: string[];
	    RotStates: Bool;
	    BldFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new BeamDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
	        this.RotStates = this.convertValues(source["RotStates"], Bool);
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
	    Lines: string[];
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
	    ShftTilt: Real;
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
	        this.Lines = source["Lines"];
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
	        this.ShftTilt = this.convertValues(source["ShftTilt"], Real);
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
	    SimProgress: number;
	    LinProgress: number;
	    LogPath: string;
	    Error: string;
	
	    static createFrom(source: any = {}) {
	        return new EvalStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.State = source["State"];
	        this.SimProgress = source["SimProgress"];
	        this.LinProgress = source["LinProgress"];
	        this.LogPath = source["LogPath"];
	        this.Error = source["Error"];
	    }
	}
	export class Evaluate {
	    ExecPath: string;
	    ExecVersion: string;
	    ExecValid: boolean;
	    NumCPUs: number;
	    Status: EvalStatus;
	
	    static createFrom(source: any = {}) {
	        return new Evaluate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ExecPath = source["ExecPath"];
	        this.ExecVersion = source["ExecVersion"];
	        this.ExecValid = source["ExecValid"];
	        this.NumCPUs = source["NumCPUs"];
	        this.Status = this.convertValues(source["Status"], EvalStatus);
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
	export class StControl {
	    Name: string;
	    Type: string;
	    Lines: string[];
	    PrescribedForcesFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new StControl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
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
	    Lines: string[];
	
	    static createFrom(source: any = {}) {
	        return new Misc(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
	    }
	}
	export class OLAF {
	    Name: string;
	    Type: string;
	    Lines: string[];
	    PrescribedCircFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new OLAF(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
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
	    Lines: string[];
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
	        this.Lines = source["Lines"];
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
	    Lines: string[];
	    PCMode: Integer;
	    VSContrl: Integer;
	    VS_RtGnSp: Real;
	    VS_RtTq: Real;
	    VS_Rgn2K: Real;
	    VS_SlPc: Real;
	    HSSBrMode: Integer;
	    YCMode: Integer;
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
	        this.Lines = source["Lines"];
	        this.PCMode = this.convertValues(source["PCMode"], Integer);
	        this.VSContrl = this.convertValues(source["VSContrl"], Integer);
	        this.VS_RtGnSp = this.convertValues(source["VS_RtGnSp"], Real);
	        this.VS_RtTq = this.convertValues(source["VS_RtTq"], Real);
	        this.VS_Rgn2K = this.convertValues(source["VS_Rgn2K"], Real);
	        this.VS_SlPc = this.convertValues(source["VS_SlPc"], Real);
	        this.HSSBrMode = this.convertValues(source["HSSBrMode"], Integer);
	        this.YCMode = this.convertValues(source["YCMode"], Integer);
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
	    Lines: string[];
	    PotFile: Path;
	
	    static createFrom(source: any = {}) {
	        return new HydroDyn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
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
	export class String {
	    Name: string;
	    Type: string;
	    Desc: string;
	    Line: number;
	    Value: string;
	
	    static createFrom(source: any = {}) {
	        return new String(source);
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
	    Lines: string[];
	    TMax: Real;
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
	    Gravity: Real;
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
	    OutFmt: String;
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
	    WrVTK: Integer;
	    VTK_type: Integer;
	
	    static createFrom(source: any = {}) {
	        return new Main(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Lines = source["Lines"];
	        this.TMax = this.convertValues(source["TMax"], Real);
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
	        this.Gravity = this.convertValues(source["Gravity"], Real);
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
	        this.OutFmt = this.convertValues(source["OutFmt"], String);
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
	        this.WrVTK = this.convertValues(source["WrVTK"], Integer);
	        this.VTK_type = this.convertValues(source["VTK_type"], Integer);
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
	
	
	export class Mode {
	    ID: number;
	    OP: number;
	    NaturalFreqHz: number;
	    DampedFreqHz: number;
	    DampingRatio: number;
	
	    static createFrom(source: any = {}) {
	        return new Mode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.OP = source["OP"];
	        this.NaturalFreqHz = source["NaturalFreqHz"];
	        this.DampedFreqHz = source["DampedFreqHz"];
	        this.DampingRatio = source["DampingRatio"];
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
	
	export class OperatingPoint {
	    ID: number;
	    Files: string[];
	    RotSpeed: number;
	    WindSpeed: number;
	    Modes: Mode[];
	
	    static createFrom(source: any = {}) {
	        return new OperatingPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Files = source["Files"];
	        this.RotSpeed = source["RotSpeed"];
	        this.WindSpeed = source["WindSpeed"];
	        this.Modes = this.convertValues(source["Modes"], Mode);
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
	
	
	
	
	
	export class Results {
	    LinDir: string;
	    HasWind: boolean;
	    OPs: OperatingPoint[];
	    LinOPs: lin.OPResults[];
	
	    static createFrom(source: any = {}) {
	        return new Results(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LinDir = source["LinDir"];
	        this.HasWind = source["HasWind"];
	        this.OPs = this.convertValues(source["OPs"], OperatingPoint);
	        this.LinOPs = this.convertValues(source["LinOPs"], lin.OPResults);
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


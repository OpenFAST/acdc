import { main } from "../wailsjs/go/models"

export type Field = main.Integer | main.Bool | main.Path | main.Paths | main.Real | main.Reals | main.String

export type ModelFile = (main.AeroDyn | main.AeroDyn14 | main.AirfoilInfo | main.BeamDyn |
    main.ElastoDyn | main.HydroDyn | main.InflowWind | main.Main | main.Misc | main.OLAF)

export interface FileOption {
    name: string;
    type: string;
    value: File;
}

export interface File {
    Name: string
    Type: string
    ID: number
    Lines: string[]
    Fields: Field[]
}

export function instanceOfField(obj: any): obj is Field {
    return typeof obj == 'object' && 'Name' in obj && 'Type' in obj;
}
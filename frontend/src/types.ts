import { main } from "../wailsjs/go/models"

export type Field = main.Integer | main.Bool | main.Path | main.Paths | main.Real | main.Reals

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
    Lines: string[]
    Fields: Field[]
}


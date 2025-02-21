// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {sqlite} from '../models';
import {main} from '../models';

export function CreateApiClient(arg1:string):Promise<void>;

export function DeleteApiClient():Promise<void>;

export function DeleteApiKey(arg1:string):Promise<void>;

export function GetAllApiKeys():Promise<Array<sqlite.ApiKey>>;

export function GetWildberriesData(arg1:string):Promise<main.ExportWbData>;

export function SaveApiKey(arg1:string,arg2:string):Promise<number>;

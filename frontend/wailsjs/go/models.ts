export namespace main {
	
	export class ExportWbData {
	    articles: string;
	    stats: string;
	
	    static createFrom(source: any = {}) {
	        return new ExportWbData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.articles = source["articles"];
	        this.stats = source["stats"];
	    }
	}

}

export namespace sqlite {
	
	export class ApiKey {
	    id: number;
	    api_key: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new ApiKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.api_key = source["api_key"];
	        this.name = source["name"];
	    }
	}

}


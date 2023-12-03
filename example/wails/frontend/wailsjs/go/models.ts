export namespace swails {
	
	export class WailsCommand {
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new WailsCommand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class WailsHTMLResponse {
	    html: string;
	
	    static createFrom(source: any = {}) {
	        return new WailsHTMLResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.html = source["html"];
	    }
	}
	export class WailsInput {
	    name: string;
	    type: string;
	    value: any;
	
	    static createFrom(source: any = {}) {
	        return new WailsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.value = source["value"];
	    }
	}

}


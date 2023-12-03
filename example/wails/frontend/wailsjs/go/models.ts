export namespace swails {
	
	export class WailsCommand {
	    name: string;
	    description: string;
	    image: string;
	    emoji: string;
	
	    static createFrom(source: any = {}) {
	        return new WailsCommand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.image = source["image"];
	        this.emoji = source["emoji"];
	    }
	}
	export class WailsHTMLResponse {
	    default: string;
	    text: string;
	    json: any;
	    table: string[][];
	    table_styles: string[][];
	
	    static createFrom(source: any = {}) {
	        return new WailsHTMLResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.default = source["default"];
	        this.text = source["text"];
	        this.json = source["json"];
	        this.table = source["table"];
	        this.table_styles = source["table_styles"];
	    }
	}
	export class WailsInput {
	    name: string;
	    type: string;
	    value: any;
	    shared: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WailsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.value = source["value"];
	        this.shared = source["shared"];
	    }
	}

}


export namespace library {
	
	export class Book {
	    filepath: string;
	    title: string;
	    format: string;
	    page?: number;
	
	    static createFrom(source: any = {}) {
	        return new Book(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filepath = source["filepath"];
	        this.title = source["title"];
	        this.format = source["format"];
	        this.page = source["page"];
	    }
	}

}

export namespace main {
	
	export class Command {
	    Name: string;
	    Run: string;
	    Key: string;
	
	    static createFrom(source: any = {}) {
	        return new Command(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Run = source["Run"];
	        this.Key = source["Key"];
	    }
	}

}


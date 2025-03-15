export namespace books {
	
	export class BookmarkInfo {
	    filename: string;
	    page: number;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new BookmarkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.page = source["page"];
	        this.title = source["title"];
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


export namespace books {
	
	export class Book {
	    id: number;
	    filepath: string;
	    title: string;
	    filesize: number;
	    // Go type: time
	    modtime: any;
	    format: string;
	
	    static createFrom(source: any = {}) {
	        return new Book(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.filepath = source["filepath"];
	        this.title = source["title"];
	        this.filesize = source["filesize"];
	        this.modtime = this.convertValues(source["modtime"], null);
	        this.format = source["format"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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


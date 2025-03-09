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


import {Lang as lang}	from "@/components/graph/utils/lang";

export class Transform {

	x:number
	y:number
	k:number
	private _element:Element;

	constructor(el:Element) {
		this._element = el;
		this.parse(el);
	}

	/**
	 * transform 파싱
	 */ 
	private parse(el:Element) {
		this._element = el;
		this.x = 0, this.y = 0, this.k = 1;

		let transform:string = el.getAttribute("transform") || "";
		if(transform) {
			let offset:number = transform.indexOf("translate");
			if(offset>=0) {
				offset+=9;
				let translate:string = transform.substring(offset).replace("(","").replace(")","").replace(","," ");
				let arr:Array<string> = translate.split(" ")
				if(arr && arr.length>=2) {
					this.x =  lang.toNumber(arr[0],0);
					this.y =  lang.toNumber(arr[1],0);
				}
			}
			offset = transform.indexOf("scale");
			if(transform && offset>=0) {
				offset+=5;
				let scale:string = transform.substring(offset).replace("(","").replace(")","").replace(","," ");
				this.k =  lang.toNumber(scale,1);
			}
		}

	}

	translate(x:number, y:number):Transform {
		return (this.x = x, this.y = y, this.transform());
	}

	scale(ratio:number):Transform {
		return (this.k = ratio, this.transform());
	}

	translateY(y:number):Transform {
		return (this.y = y, this.transform());
	}
	translateX(x:number):Transform {
		return (this.x = x, this.transform());
	}

	toString():string {
		let attr:string = `translate(${this.x} ${this.y})`;
		return this.k>0 ? `${attr} scale(${this.k})`: attr;
	}

	private transform():Transform {
		this._element ? (this._element.setAttribute("transform",	this.toString()), this): this;
		return this;
	}


	static instance(el:Element):Transform {
		return new Transform(el);
	}

}

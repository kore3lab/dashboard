
export class Lang {
	
	public static toNumber(n:any, d:number):number {
		if(typeof(n) == "string") n = n.replace(/\px/g,'')
		return Lang.isNumber(n)? Number.parseFloat(n): (d?d:0);
	}

	public static isNumber(n:any):boolean{
		return !isNaN(parseFloat(n)) && isFinite(n);
	}

	/**
	 * 객체 복사하기 (덮어쓰기. 이전 데이터 유지)
	 * @param dest to (source 가 object 인 경우만 값을 가지고 있음, 기존객체 값 유지위해서)
	 * @param source from
	 */
	public static merge(dest:any, source:any):any {

		if (null == source || "object" != typeof source) return source;

		let copy:any;
	
		// Handle Date
		if (source instanceof Date) {
			copy = new Date();
			copy.setTime(source.getTime());
			return copy;
		}
	
		// Handle Array
		if (source instanceof Array) {
			copy = [];
			for (let i = 0, len = source.length; i < len; i++) {
				copy[i] = this.merge(null, source[i]);
			}
			return copy;
		}
	
		// Handle Object
		if (dest instanceof Object || source instanceof Object) {
			copy = {};
			for (let attr in dest) {
				if (dest.hasOwnProperty(attr)) copy[attr] = this.merge((dest[attr] instanceof Object)?copy[attr]: null, dest[attr]);
			}

			for (let attr in source) {
				if (source.hasOwnProperty(attr)) copy[attr] = this.merge((source[attr] instanceof Object)?copy[attr]:null, source[attr]);
			}

			return copy;
		}
	
		return source;
	}



}


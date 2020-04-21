import {Transform}	from "./transform";
import {Bounds}		from "../model/models";

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

export class UI {

	/**
	 * 주어진 element 의 부모 element 높이 맞추어 수직정렬
	 * 
	 */	
	public static alignVertical(el:SVGElement) {
		
		if (el ==null || el.parentElement == null) return 

		const rect:ClientRect = el.getBoundingClientRect();
		const rectParent:ClientRect = el.parentElement.getBoundingClientRect();

		if(rect.height < rectParent.height) {
			Transform.instance(el).shiftY((rectParent.height-rect.height)/2);
		}

	}

	/**
	 * 주어진 element 의 부모 element 너비 맞추어 수평정렬
	 * 
	 */	
	public static alignHorizonal(el:SVGElement) {

		if (el ==null || el.parentElement == null) return 

		const rect:ClientRect = el.getBoundingClientRect();
		const rectParent:ClientRect = el.parentElement.getBoundingClientRect();

		if(rect.width < rectParent.width) {
			Transform.instance(el).shiftX((rectParent.width-rect.width)/2);
		}

	}

	/**
	 * 같은 형재 element 들을 가운데 수평 정렬
	 * 
	 */	
	public static alignHorizonals(els:Array<SVGElement>) { 

		if(els.length < 1) return;
		let maxWidth:number = 0;

		els.forEach( (el:SVGElement) => {
			const rect:ClientRect = el.getBoundingClientRect();
			maxWidth = Math.max(rect.width+rect.left, maxWidth)
		});

		els.forEach( (el:SVGElement) => {
			const rect:ClientRect = el.getBoundingClientRect();
			if(rect.width < maxWidth) {
				Transform.instance(el).shiftX((maxWidth-rect.width)/2);
			}
		});


	}


	/**
	 * 스크롤 가능한 레이어 추가
	 * 
	 * @param X 추가할 레이어의 X 위치
	 * @param Y 추가할 레이어의 Y 위치
	 * @param bounds 기준이 되는 bounds (스크롤 여부를 결정(계산)할 때 사용)
	 * @param parentEl 추가할 레이어 
	 * @param func 레이어의 내용 함수
	 * @param args 레이어의 내용 함수의 파라메터
	 */
	public static appendScrollableLayer(X:number, Y:number, bounds:Bounds, parentEl:d3.Selection<SVGGElement,any,SVGElement,any>, func: (selection: d3.Selection<SVGElement, any, SVGElement, any>, ...args: any) => void , ...args: any[]) {

		const margin:number = 10;	//마진 기준 (top,left, right, bottom 동일 처리)

		// 스크롤을 위해서 "div" 사용을 위해서 "foreignObject" 엘리먼트 활용
		let scrollEl:d3.Selection<SVGForeignObjectElement,any,SVGElement,any> = parentEl.append("foreignObject")
			.attr("x",X)
			.attr("y",Y)
			.html(`<div xmlns="http://www.w3.org/1999/xhtml" style="height: 100%;padding:${margin}px;"></div>`)

		// div 엘리먼트에 svg 추가
		let svg:d3.Selection<SVGSVGElement,any,SVGElement,any> = scrollEl.select("div").append("svg")
		
		// div 엘리먼트에 g.outline 추가
		let outlineEl:d3.Selection<SVGGElement,any,SVGElement,any> = svg.append("g").attr("class","outline")

		// outline 엘리먼트에 파라메터로 받은 render 실행
		outlineEl.call(func, args);

			
		// 그려진 outline Bounds
		const rect:ClientRect = outlineEl.node()!.getBoundingClientRect();

		// 범례 백그라운드
		outlineEl.insert("rect", "g.group:first-child")
			.attr("class", "background")
			.attr("width", rect.width)
			.attr("height", rect.height)

		// foreignObject / svg 엘리먼트 너비
		scrollEl.attr("width", rect.width + (margin*2));
		svg.attr("width", rect.width + (margin*2))
		
		// 스크롤 위한 높이 정의 (내용이 더 크면 스크롤이 생기도록함)
		const outH:number = bounds.height - (Y+margin)-margin;								//outer 높이 : bottom 마진 + outline bottom 마진 반영(x2)
		const inH:number = outlineEl!.node()!.getBoundingClientRect().height + (margin*2);	//innert 높이

		scrollEl.attr("height", (inH>outH) ? outH:inH);										//내용 높이(inenr)가 더 크면 outer , 내용 높으가 더 작다면 inner 로 줄임
		if(inH>outH) scrollEl.select("div").style("overflow-y", "scroll");
		svg.attr("height", inH );

	}



}

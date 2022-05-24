import * as d3		from "d3";
import {Transform}	from "@/components/graph/utils/transform";
import {Lang}		from "@/components/graph/utils/lang";
import { select } from "d3";

export class WH {
    height: number;
    width: number;
}

export class Bounds extends WH {

    x: number;
    y: number;

	constructor(el:d3.Selection<SVGGElement, any, Element, any>|HTMLElement) {
		super();
		if (el instanceof HTMLElement) {
			const bounds:DOMRect = el.getBoundingClientRect();	//absolute-position
			const selection = d3.select(el);
			this.x	= bounds.x + Lang.toNumber(selection.style("padding-left"),0);
			this.y	= bounds.y + Lang.toNumber(selection.style("padding-top"),0);
			this.width	= bounds.right - Lang.toNumber(selection.style("padding-right"),0) - this.x;
			this.height	= bounds.bottom - Lang.toNumber(selection.style("padding-bottom"),0) - this.y;
		} else {
			const bounds:DOMRect = el.node()!.getBBox();	//relative-position
			const transform:Transform = Transform.instance(el.node()!)
			this.x	= transform.x;
			this.y	= transform.y;
			this.width	= bounds.width;
			this.height	= bounds.height;
		}

	}

}

export class UI {


	/**
	 * align
	 * 
	 * @param el target element
	 * @param horizontal
	 * @param vertical
	 * @param margin margin 반영
	 */
	public static align(el:SVGElement, horizontal:"none"|"left"|"right"|"center", vertical:"none"|"top"|"bottom"|"middle",  margin?: {left?:number, top?:number, right?:number,bottom?:number}) { 

		if(!el || !el.parentElement) return
		const outBounds:DOMRect = el.parentElement.getBoundingClientRect();
		const bounds:DOMRect = el.getBoundingClientRect();

		let X:number=0, Y:number=0;
		if(!margin) margin = {left:0, top:0, right:0, bottom:0}

		if(horizontal == "right") {
			X = outBounds.width - bounds.width - ( margin.right ? margin.right: 0);
		} else if(horizontal == "center") {
			X = (outBounds.width - bounds.width - (margin.right ? margin.right: 0))/2;

		} else {
			X = (margin && margin.left) ? margin.left: 0;
		}

		if(vertical == "bottom") {
			Y = outBounds.height - bounds.height - (margin.bottom ? margin.bottom: 0) ;
		} else if(vertical == "middle") {
			Y = (outBounds.height - bounds.height - (margin.bottom ? margin.bottom: 0))/2;
		} else {
			Y = margin.top ? margin.top: 0;
		}

		const k:number = Transform.instance(el.parentElement).k;	//calcuate parent elements's tranfrom-scale (ratio) 
		Transform.instance(el).translate(X/k,Y/k)
	}

	

	public static align2(selection: d3.Selection<SVGGElement, any, SVGElement, any>, horizontal:"none"|"left"|"right"|"center", vertical:"none"|"top"|"bottom"|"middle",  margin?: {left?:number, top?:number, right?:number,bottom?:number}): d3.ZoomTransform{ 

		const el:SVGElement = selection.node()!
		let transform:d3.ZoomTransform = d3.zoomTransform(el);
		if(!el || !el.parentElement) return transform

		const outline:DOMRect = el.parentElement.getBoundingClientRect();
		const inline:DOMRect = el.getBoundingClientRect();


		let X = outline.width > inline.width ? inline.x + (outline.width - inline.width)/2 : 0;
		let Y = outline.height > inline.height ? -inline.y + (outline.height - inline.height)/2: -inline.y;

		const k:number = Transform.instance(el.parentElement).k;	//calcuate parent elements's tranfrom-scale (ratio) 
		return transform.translate(X/k,Y/k);
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
	public static appendScrollableLayer(X:number, Y:number, bounds:Bounds, parentEl:d3.Selection<SVGGElement,any,SVGElement,any>, func: (selection: d3.Selection<SVGGElement, any, SVGElement, any>, ...args: any) => void , ...args: any[]) {

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
		const rect:DOMRect = outlineEl.node()!.getBoundingClientRect();

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

	/**
	 * "text" element 말 줄임 
	 *  - width 초과되면 "..." 말줄임
	 *
	 * @param el "text" element
	 * @param width 최대 너비
	 */
	public static ellipsisText(el:SVGTextElement, width:number): void {
		width -= el.x.baseVal[0].value;	//x 값 빼기
		if(el.getComputedTextLength() > width) {
			const text = `${el.textContent}`;
			const chars = text.split("");
			let len:number = chars.length - 3;
			while(len > 1) {
				len--;
				el.textContent = `${text.substring(0, len)}...`
				if(el.getComputedTextLength() < width) break;
			}
		}
	}

	public static appendBox(parentEl:d3.Selection<SVGGElement, any, SVGElement,any>, 
		render: (selection: d3.Selection<SVGGElement, any, SVGElement, any>, ...args: any[]) => void, 
		width?:number,
		padding?:{top:number,left:number, right:number, bottom:number},
		bg?:{fill:string, opacity?:number},
		border?:{width:number, color?:string, dash?:string}): d3.Selection<SVGGElement, any, SVGElement,any> {

		const boxWrap = parentEl.append("g").attr("class","boxWrap");

		//rendering
		const box = boxWrap.append("g").attr("class","box");
		if(render) box.call(render)

		// after rendering > calcuate border(background) bounds
		let bounds:DOMRect = boxWrap.node()!.getBBox();
		const bottom:number = bounds.y + bounds.height + (padding?padding.top+padding.bottom:0);
		const right:number = bounds.x + (width?width:bounds.width); //stroke-width 반영

		// box (insert before g.box)
		const background = boxWrap.insert("path","g.box")
			.attr("class","background")
			.attr("d",`M${bounds.x},${bounds.y} L${right},${bounds.y} L${right},${bottom} L${bounds.x},${bottom} L${bounds.x},${bounds.y}`)

		if(bg) {
			background.attr("fill",bg.fill)
			if(bg.fill != "none" && bg.opacity) background.attr("fill-opacity",bg.opacity)
		}

		if(border) {
			background.attr("stroke","black").attr("stroke-width",border.width)
			if (border.color)  background.attr("stroke", border.color)
			if (border.dash) background.attr("stroke-dasharray", border.dash)
		}
			

		if(padding) Transform.instance(box.node()!).translate(padding.left,padding.top)
		return box

	}

}

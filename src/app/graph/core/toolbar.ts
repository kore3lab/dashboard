import * as d3		from "d3";
import {UI}			from "../utils/lang";
import {GraphBase}	from "./graph.base";
import "./toolbar.css";

export class LegendModel {
	header:string
	rows:Array<{label:string, ico:string}> = []
}



/**
 * 범례
 */
export class Toolbar {

	/**
	 * 툴바 랜더링
	 * 		-  범례 랜더링 포함
	 * 
	 */
	public static render(svgEl:d3.Selection<SVGSVGElement,any,SVGElement,any>, owner:GraphBase, legends:Array<LegendModel>) {

		const margin:number = 10;
		let Y:number = margin;

		if(svgEl.select("g.toolbar").size()>0) return;

		let toolbarEl:d3.Selection<SVGGElement,any,SVGElement,any> = svgEl.append("g")
			.attr("class","toolbar")
			.attr("transform",`translate (${margin}, ${Y})`)

		// zoom - in 
		toolbarEl.append("g")
			.attr("id","ac_btn_zoomin")
			.attr("class","button")
			.html('<rect></rect><g class="ico"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/><path fill="none" d="M0 0h24v24H0V0z"/><path d="M12 10h-2v2H9v-2H7V9h2V7h1v2h2v1z"/></svg></g>')
			.on("click", () => {
				owner.zoomRatio(1.1);
			});
		
		// zoom - out
		toolbarEl.append("g")
			.attr("id","ac_btn_zoomout")
			.attr("class","button")
			.html('<rect></rect><g class="ico"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" d="M0 0h24v24H0V0z"/><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14zM7 9h5v1H7z"/></svg></g>')
			.on("click", () => {
				owner.zoomRatio(0.9);
			});

		// zoom - 맞춤
		toolbarEl.append("g")
			.attr("id","ac_btn_zoomfit")
			.attr("class","button")
			.html('<rect></rect><g class="ico"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="24" height="24" viewBox="0 0 24 24"><defs><path id="a" d="M0 0h24v24H0z"/></defs><clipPath id="b"><use xlink:href="#a" overflow="visible"/></clipPath><path clip-path="url(#b)" d="M15 3l2.3 2.3-2.89 2.87 1.42 1.42L18.7 6.7 21 9V3zM3 9l2.3-2.3 2.87 2.89 1.42-1.42L6.7 5.3 9 3H3zm6 12l-2.3-2.3 2.89-2.87-1.42-1.42L5.3 17.3 3 15v6zm12-6l-2.3 2.3-2.87-2.89-1.42 1.42 2.89 2.87L15 21h6z"/><path clip-path="url(#b)" fill="none" d="M0 0h24v24H0z"/></svg></g>')
			.on("click", () => {
				owner.zoom();
			});

		// 범례 버튼
		toolbarEl.append("g")
			.attr("id","ac_btn_legend")
			.attr("class","button")
			.on("click", () => {
				let g = d3.select("g.legend");
				g.attr("visibility", g.attr("visibility")=="visible"?"hidden":"visible");
			})
			.html('<rect></rect><text>범례</text>')

		// 범례 들어갈 틀 만들기 : g.legend > foreignObject > div (스크롤) > svg > g.outline
		Y += toolbarEl.node().getBoundingClientRect().height + margin;	//Y 값 - 툴바 버튼 높이 반영

		// 범레 g.legend 추가
		let legendEl:d3.Selection<SVGGElement,any,SVGElement,any> = svgEl.append("g")
			.attr("class","legend")
			.attr("visibility", "visible")


		// g.legend  에 스크롤 가능한 레이어 추가
		UI.appendScrollableLayer(margin, Y, owner.bounds(), legendEl, Toolbar.renderLegends, legends);


		// 닫기 버튼
		legendEl.selectAll("g.outline").append("g")
			.attr("class","button")
			.attr("height", "24").attr("width", "24")
			.attr("xlink:href", "#ac_ic_close")
			.on("click",() => {
				d3.select("g.legend").attr("visibility","hidden");
			})
			.attr("transform", `translate (${legendEl.node().getBoundingClientRect().width -(margin*2+24)}, 0)`)
			.html(`<rect></rect><g class="ico"><use width="24" height="24" href="#ac_ic_close"></use></g>`)
		

	}

	/**
	 * 범례 랜더링
	 * 
	 */
	private static renderLegends(outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>, args: any) {

		let legends:Array<LegendModel> = args[0]

		//  범례 데이터로 y 값 지정 작업
		let groupY:number = 0;		// 그룹별 y 값
		const marginH:number = 30, rowH:number =25	//group 간 마진, 아이템별 높이ㄴ
		legends.forEach((d:LegendModel) => {
			d["y"] = groupY;
			groupY += (d.rows.length+1) * rowH + marginH;
			let y:number = 0;
			d.rows.forEach((r:any) => {
				r["y"] = (y+=rowH);
			});
		});

		// 범례 그리기
		// g.outline엘리먼트에 범례 그룹 추가
		let group:d3.Selection<SVGGElement,any,SVGElement,any> = outlineEl.selectAll("g")
			.data(legends)
			.enter()
			.append("g")
			.attr("transform", (d:any)=>`translate (0, ${d.y})`)
			.attr("class","group")

		// 범례 그룹에 제목 추가
		group.append("text")
			.text((d:any)=> d.header)

		// 범례 그룹별 아이템 리스트 추가
		group
			.selectAll("g.rows")
			.data(d=>d.rows)
			.enter()
			.append("g")
			.attr("class", "row")
			.html((d:any)=> {
				return `${d.ico}<text x="50">${d.label}</text>`;
			})
			.attr("transform", (d:any)=> {return `translate (0, ${d.y})`})


	}
}

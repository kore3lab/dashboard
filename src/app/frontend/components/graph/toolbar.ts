import * as d3		from "d3";
import {GraphBase}	from "@/components/graph/graph.base";
import "@/components/graph/toolbar.css";

/**
 * Toolbar
 */
export class Toolbar {

	/**
	 * Redering toolbar
	 */
	public static render(owner:GraphBase):d3.Selection<SVGGElement,any,SVGElement,any> {

		const svgEl:d3.Selection<SVGSVGElement,any,SVGElement,any> = owner.svg;

		// root
		let toolbarEl:d3.Selection<SVGGElement,any,SVGElement,any> = svgEl.select("g.toolbar");
		if(toolbarEl.size() == 0) toolbarEl = svgEl.append("g").attr("class","toolbar")
			
		// button - zoom in 
		toolbarEl.append("g")
			.attr("id","ac_btn_zoomin")
			.attr("class","button")
			.html('<rect></rect><g class="ico"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/><path fill="none" d="M0 0h24v24H0V0z"/><path d="M12 10h-2v2H9v-2H7V9h2V7h1v2h2v1z"/></svg></g>')
			.on("click", () => {
				owner.zoomRatio(1.1);
			});
		
		// button - zoom out
		toolbarEl.append("g")
			.attr("id","ac_btn_zoomout")
			.attr("class","button")
			.html('<rect></rect><g class="ico"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" d="M0 0h24v24H0V0z"/><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14zM7 9h5v1H7z"/></svg></g>')
			.on("click", () => {
				owner.zoomRatio(0.9);
			});

		// button - zoom fit
		toolbarEl.append("g")
			.attr("id","ac_btn_zoomfit")
			.attr("class","button")
			.html('<rect></rect><g class="ico"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="24" height="24" viewBox="0 0 24 24"><defs><path id="a" d="M0 0h24v24H0z"/></defs><clipPath id="b"><use xlink:href="#a" overflow="visible"/></clipPath><path clip-path="url(#b)" d="M15 3l2.3 2.3-2.89 2.87 1.42 1.42L18.7 6.7 21 9V3zM3 9l2.3-2.3 2.87 2.89 1.42-1.42L6.7 5.3 9 3H3zm6 12l-2.3-2.3 2.89-2.87-1.42-1.42L5.3 17.3 3 15v6zm12-6l-2.3 2.3-2.87-2.89-1.42 1.42 2.89 2.87L15 21h6z"/><path clip-path="url(#b)" fill="none" d="M0 0h24v24H0z"/></svg></g>')
			.on("click", () => {
				owner.zoom();
			});

		return toolbarEl;

	}

}

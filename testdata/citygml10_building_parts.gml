<?xml version="1.0" encoding="UTF-8"?>
<!-- Synthetic CityGML 1.0 LoD2 fixture shaped like an AdV/LGLN tile: a
     building made of two BuildingParts with their own measuredHeight,
     lod2Solid (xlink:href members only) and boundedBy surfaces, and a plain
     building whose GroundSurface has two polygons, the larger with a hole. -->
<core:CityModel xmlns:core="http://www.opengis.net/citygml/1.0"
                xmlns:bldg="http://www.opengis.net/citygml/building/1.0"
                xmlns:gen="http://www.opengis.net/citygml/generics/1.0"
                xmlns:app="http://www.opengis.net/citygml/appearance/1.0"
                xmlns:grp="http://www.opengis.net/citygml/cityobjectgroup/1.0"
                xmlns:gml="http://www.opengis.net/gml"
                xmlns:xlink="http://www.w3.org/1999/xlink">
  <gml:name>synthetic_lod2_citygml10</gml:name>
  <gml:boundedBy>
    <gml:Envelope srsName="urn:adv:crs:ETRS89_UTM32*DE_DHHN2016_NH">
      <gml:lowerCorner srsDimension="3">550000 5803000 50</gml:lowerCorner>
      <gml:upperCorner srsDimension="3">550100 5803100 62</gml:upperCorner>
    </gml:Envelope>
  </gml:boundedBy>
  <core:cityObjectMember>
    <bldg:Building gml:id="BLDG_PARTS">
      <core:creationDate>2024-01-07</core:creationDate>
      <gen:stringAttribute name="Gemeindeschluessel">
        <gen:value>00000000</gen:value>
      </gen:stringAttribute>
      <bldg:function>31001_1000</bldg:function>
      <bldg:consistsOfBuildingPart>
        <bldg:BuildingPart gml:id="PART_A">
          <gen:stringAttribute name="DachtypName">
            <gen:value>PolyFlatRoof</gen:value>
          </gen:stringAttribute>
          <bldg:function>31001_1000</bldg:function>
          <bldg:roofType>1000</bldg:roofType>
          <bldg:measuredHeight uom="urn:adv:uom:m">12.0</bldg:measuredHeight>
          <bldg:lod2Solid>
            <gml:Solid>
              <gml:exterior>
                <gml:CompositeSurface>
                  <gml:surfaceMember xlink:href="#PA_G"/>
                  <gml:surfaceMember xlink:href="#PA_R"/>
                  <gml:surfaceMember xlink:href="#PA_W1"/>
                  <gml:surfaceMember xlink:href="#PA_W_MISSING"/>
                </gml:CompositeSurface>
              </gml:exterior>
            </gml:Solid>
          </bldg:lod2Solid>
          <bldg:boundedBy>
            <bldg:GroundSurface gml:id="PA_GS">
              <bldg:lod2MultiSurface>
                <gml:MultiSurface>
                  <gml:surfaceMember>
                    <gml:Polygon gml:id="PA_G">
                      <gml:exterior>
                        <gml:LinearRing>
                          <gml:posList srsDimension="3">550000 5803000 50 550000 5803010 50 550010 5803010 50 550010 5803000 50 550000 5803000 50</gml:posList>
                        </gml:LinearRing>
                      </gml:exterior>
                    </gml:Polygon>
                  </gml:surfaceMember>
                </gml:MultiSurface>
              </bldg:lod2MultiSurface>
            </bldg:GroundSurface>
          </bldg:boundedBy>
          <bldg:boundedBy>
            <bldg:RoofSurface gml:id="PA_RS">
              <bldg:lod2MultiSurface>
                <gml:MultiSurface>
                  <gml:surfaceMember>
                    <gml:Polygon gml:id="PA_R">
                      <gml:exterior>
                        <gml:LinearRing>
                          <gml:posList srsDimension="3">550000 5803000 62 550010 5803000 62 550010 5803010 62 550000 5803010 62 550000 5803000 62</gml:posList>
                        </gml:LinearRing>
                      </gml:exterior>
                    </gml:Polygon>
                  </gml:surfaceMember>
                </gml:MultiSurface>
              </bldg:lod2MultiSurface>
            </bldg:RoofSurface>
          </bldg:boundedBy>
          <bldg:boundedBy>
            <bldg:WallSurface gml:id="PA_WS">
              <bldg:lod2MultiSurface>
                <gml:MultiSurface>
                  <gml:surfaceMember>
                    <gml:Polygon gml:id="PA_W1">
                      <gml:exterior>
                        <gml:LinearRing>
                          <gml:posList srsDimension="3">550000 5803000 50 550010 5803000 50 550010 5803000 62 550000 5803000 62 550000 5803000 50</gml:posList>
                        </gml:LinearRing>
                      </gml:exterior>
                    </gml:Polygon>
                  </gml:surfaceMember>
                </gml:MultiSurface>
              </bldg:lod2MultiSurface>
            </bldg:WallSurface>
          </bldg:boundedBy>
        </bldg:BuildingPart>
      </bldg:consistsOfBuildingPart>
      <bldg:consistsOfBuildingPart>
        <bldg:BuildingPart gml:id="PART_B">
          <bldg:function>31001_2000</bldg:function>
          <bldg:measuredHeight uom="urn:adv:uom:m">6.5</bldg:measuredHeight>
          <bldg:lod2Solid>
            <gml:Solid>
              <gml:exterior>
                <gml:CompositeSurface>
                  <gml:surfaceMember xlink:href="#PB_G"/>
                  <gml:surfaceMember xlink:href="#PB_R"/>
                  <gml:surfaceMember xlink:href="#PB_W1"/>
                  <gml:surfaceMember xlink:href="#PB_W_MISSING"/>
                </gml:CompositeSurface>
              </gml:exterior>
            </gml:Solid>
          </bldg:lod2Solid>
          <bldg:boundedBy>
            <bldg:GroundSurface gml:id="PB_GS">
              <bldg:lod2MultiSurface>
                <gml:MultiSurface>
                  <gml:surfaceMember>
                    <gml:Polygon gml:id="PB_G">
                      <gml:exterior>
                        <gml:LinearRing>
                          <gml:posList srsDimension="3">550010 5803000 50 550010 5803008 50 550016 5803008 50 550016 5803000 50 550010 5803000 50</gml:posList>
                        </gml:LinearRing>
                      </gml:exterior>
                    </gml:Polygon>
                  </gml:surfaceMember>
                </gml:MultiSurface>
              </bldg:lod2MultiSurface>
            </bldg:GroundSurface>
          </bldg:boundedBy>
          <bldg:boundedBy>
            <bldg:RoofSurface gml:id="PB_RS">
              <bldg:lod2MultiSurface>
                <gml:MultiSurface>
                  <gml:surfaceMember>
                    <gml:Polygon gml:id="PB_R">
                      <gml:exterior>
                        <gml:LinearRing>
                          <gml:posList srsDimension="3">550010 5803000 56.5 550016 5803000 56.5 550016 5803008 56.5 550010 5803008 56.5 550010 5803000 56.5</gml:posList>
                        </gml:LinearRing>
                      </gml:exterior>
                    </gml:Polygon>
                  </gml:surfaceMember>
                </gml:MultiSurface>
              </bldg:lod2MultiSurface>
            </bldg:RoofSurface>
          </bldg:boundedBy>
          <bldg:boundedBy>
            <bldg:WallSurface gml:id="PB_WS">
              <bldg:lod2MultiSurface>
                <gml:MultiSurface>
                  <gml:surfaceMember>
                    <gml:Polygon gml:id="PB_W1">
                      <gml:exterior>
                        <gml:LinearRing>
                          <gml:posList srsDimension="3">550010 5803000 50 550016 5803000 50 550016 5803000 56.5 550010 5803000 56.5 550010 5803000 50</gml:posList>
                        </gml:LinearRing>
                      </gml:exterior>
                    </gml:Polygon>
                  </gml:surfaceMember>
                </gml:MultiSurface>
              </bldg:lod2MultiSurface>
            </bldg:WallSurface>
          </bldg:boundedBy>
        </bldg:BuildingPart>
      </bldg:consistsOfBuildingPart>
      <bldg:address>
        <core:Address>
          <core:xalAddress/>
        </core:Address>
      </bldg:address>
    </bldg:Building>
  </core:cityObjectMember>
  <core:cityObjectMember>
    <bldg:Building gml:id="BLDG_PLAIN">
      <bldg:function>31001_1010</bldg:function>
      <bldg:measuredHeight uom="urn:adv:uom:m">8</bldg:measuredHeight>
      <bldg:lod2Solid>
        <gml:Solid>
          <gml:exterior>
            <gml:CompositeSurface>
              <gml:surfaceMember xlink:href="#PL_G1"/>
              <gml:surfaceMember xlink:href="#PL_G2"/>
              <gml:surfaceMember xlink:href="#PL_R"/>
            </gml:CompositeSurface>
          </gml:exterior>
        </gml:Solid>
      </bldg:lod2Solid>
      <bldg:lod2TerrainIntersection>
        <gml:MultiCurve>
          <gml:curveMember>
            <gml:LineString>
              <gml:posList srsDimension="3">550040 5803040 50 550050 5803040 50</gml:posList>
            </gml:LineString>
          </gml:curveMember>
        </gml:MultiCurve>
      </bldg:lod2TerrainIntersection>
      <bldg:boundedBy>
        <bldg:GroundSurface gml:id="PL_GS">
          <bldg:lod2MultiSurface>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon gml:id="PL_G1">
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList srsDimension="3">550030 5803040 50 550030 5803042 50 550032 5803042 50 550032 5803040 50 550030 5803040 50</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
              <gml:surfaceMember>
                <gml:Polygon gml:id="PL_G2">
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList srsDimension="3">550040 5803040 50 550040 5803050 50 550050 5803050 50 550050 5803040 50 550040 5803040 50</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                  <gml:interior>
                    <gml:LinearRing>
                      <gml:posList srsDimension="3">550044 5803044 50 550046 5803044 50 550046 5803046 50 550044 5803046 50 550044 5803044 50</gml:posList>
                    </gml:LinearRing>
                  </gml:interior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </bldg:lod2MultiSurface>
        </bldg:GroundSurface>
      </bldg:boundedBy>
      <bldg:boundedBy>
        <bldg:RoofSurface gml:id="PL_RS">
          <bldg:lod2MultiSurface>
            <gml:MultiSurface>
              <gml:surfaceMember>
                <gml:Polygon gml:id="PL_R">
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList srsDimension="3">550040 5803040 58 550050 5803040 58 550050 5803050 58 550040 5803050 58 550040 5803040 58</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:MultiSurface>
          </bldg:lod2MultiSurface>
        </bldg:RoofSurface>
      </bldg:boundedBy>
    </bldg:Building>
  </core:cityObjectMember>
</core:CityModel>

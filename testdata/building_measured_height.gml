<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="EPSG:25832">
      <gml:lowerCorner>500000 5700000 0</gml:lowerCorner>
      <gml:upperCorner>500020 5700020 15</gml:upperCorner>
    </gml:Envelope>
  </gml:boundedBy>
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_MH_001">
      <bldg:class>1000</bldg:class>
      <bldg:function>1010</bldg:function>
      <bldg:usage>residential</bldg:usage>
      <bldg:measuredHeight>12.5</bldg:measuredHeight>
      <bldg:lod1Solid>
        <gml:Solid>
          <gml:exterior>
            <gml:CompositeSurface>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>500000 5700000 0 500020 5700000 0 500020 5700020 0 500000 5700020 0 500000 5700000 0</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
              <gml:surfaceMember>
                <gml:Polygon>
                  <gml:exterior>
                    <gml:LinearRing>
                      <gml:posList>500000 5700000 12.5 500020 5700000 12.5 500020 5700020 12.5 500000 5700020 12.5 500000 5700000 12.5</gml:posList>
                    </gml:LinearRing>
                  </gml:exterior>
                </gml:Polygon>
              </gml:surfaceMember>
            </gml:CompositeSurface>
          </gml:exterior>
        </gml:Solid>
      </bldg:lod1Solid>
    </bldg:Building>
  </cityObjectMember>
</CityModel>

<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="EPSG:25832">
      <gml:lowerCorner>500000 5700000 100</gml:lowerCorner>
      <gml:upperCorner>500020 5700015 108</gml:upperCorner>
    </gml:Envelope>
  </gml:boundedBy>
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_ZE_001">
      <bldg:lod1MultiSurface>
        <gml:MultiSurface>
          <gml:surfaceMember>
            <gml:Polygon>
              <gml:exterior>
                <gml:LinearRing>
                  <gml:posList>500000 5700000 100 500020 5700000 100 500020 5700015 100 500000 5700015 100 500000 5700000 100</gml:posList>
                </gml:LinearRing>
              </gml:exterior>
            </gml:Polygon>
          </gml:surfaceMember>
          <gml:surfaceMember>
            <gml:Polygon>
              <gml:exterior>
                <gml:LinearRing>
                  <gml:posList>500000 5700000 108 500020 5700000 108 500020 5700015 108 500000 5700015 108 500000 5700000 108</gml:posList>
                </gml:LinearRing>
              </gml:exterior>
            </gml:Polygon>
          </gml:surfaceMember>
        </gml:MultiSurface>
      </bldg:lod1MultiSurface>
    </bldg:Building>
  </cityObjectMember>
</CityModel>

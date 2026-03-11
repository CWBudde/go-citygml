<?xml version="1.0" encoding="UTF-8"?>
<CityModel xmlns="http://www.opengis.net/citygml/2.0"
           xmlns:gml="http://www.opengis.net/gml"
           xmlns:bldg="http://www.opengis.net/citygml/building/2.0">
  <gml:boundedBy>
    <gml:Envelope srsName="EPSG:25832">
      <gml:lowerCorner>500000 5700000 0</gml:lowerCorner>
      <gml:upperCorner>500100 5700100 30</gml:upperCorner>
    </gml:Envelope>
  </gml:boundedBy>
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_A">
      <bldg:measuredHeight>5.0</bldg:measuredHeight>
    </bldg:Building>
  </cityObjectMember>
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_B">
      <bldg:measuredHeight>10.0</bldg:measuredHeight>
      <bldg:class>1120</bldg:class>
    </bldg:Building>
  </cityObjectMember>
  <cityObjectMember>
    <bldg:Building gml:id="BLDG_C">
      <bldg:measuredHeight>25.0</bldg:measuredHeight>
      <bldg:function>commercial</bldg:function>
    </bldg:Building>
  </cityObjectMember>
</CityModel>

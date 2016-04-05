Siemens MRI Raw Data Tools
---------------------------

A couple of simple tools to create a (somewhat) unique signature of `*.dat` files and DICOM files from siemens MRI systems

To install:

    go get github.com/hansenms/siemensraw/...


To get a signature from a datfile:

    datsignature <DATFILENAME>

To get a signature from a DICOM:

   dicomsignature <DICOMFILE>


If signatures from DICOM and `*.dat` file match, there is a pretty good chance that the images were generated from the `*.dat` file. 
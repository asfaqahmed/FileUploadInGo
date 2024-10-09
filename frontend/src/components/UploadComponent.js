import React, { useState } from 'react';
import axios from 'axios';

const UploadComponent = () => {
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState("");
  const [uploadedFileID, setUploadedFileID] = useState(null);

  const onFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const onFileUpload = async (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("file", file);

    const API_URL = process.env.REACT_APP_API_URL;

    try {
      const res = await axios.post(`${API_URL}/upload`, formData);
      setMessage(res.data.message);
      setUploadedFileID(res.data.fileID);  // Store the file ID for retrieval
    } catch (err) {
      console.error(err);
      setMessage("File upload failed!");
    }
  };

  const getUploadedFileURL = () => {
    const API_URL = process.env.REACT_APP_API_URL;
    return `${API_URL}/files/${uploadedFileID}`;
  };

  return (
    <div>
      <h2>File Upload</h2>
      <form onSubmit={onFileUpload}>
        <input type="file" onChange={onFileChange} />
        <button type="submit">Upload</button>
      </form>
      <p>{message}</p>
      {uploadedFileID && (
        <div>
          <h3>Uploaded File:</h3>
          <img src={getUploadedFileURL()} alt="Uploaded" />
        </div>
      )}
    </div>
  );
};

export default UploadComponent;

import React, { useState } from 'react';
import axios from 'axios';

const UploadComponent = () => {
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState("");

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
    } catch (err) {
      console.error(err);
      setMessage("File upload failed!");
    }
  };

  return (
    <div>
      <h2>File Upload</h2>
      <form onSubmit={onFileUpload}>
        <input type="file" onChange={onFileChange} />
        <button type="submit">Upload</button>
      </form>
      <p>{message}</p>
    </div>
  );
};

export default UploadComponent;
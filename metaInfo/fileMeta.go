package metaInfo

import mydb "myFileServer/db"

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1, FileName, FileLocation, UpdateTime string
	FileSize                                     int64
}

var fileMetaMap map[string]FileMeta

func init() {
	fileMetaMap = make(map[string]FileMeta)
}

// UpdateFileMetaMap 新增/更新文件元信息
func UpdateFileMetaMap(newMeta FileMeta) {
	fileMetaMap[newMeta.FileSha1] = newMeta
}

// UpdateFileMetaDB 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(newMeta FileMeta) bool {
	return mydb.OnFileUploadFinished(newMeta.FileSha1, newMeta.FileName, newMeta.FileLocation, newMeta.FileSize)
}

// GetFileMeta 通过sha1值获取某个文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetaMap[fileSha1]
}

// GetFileMetaDB 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tFile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, nil
	}
	fMeta := FileMeta{
		FileSha1:     tFile.FileHash,
		FileName:     tFile.FileName.String,
		FileLocation: tFile.FileAddr.String,
		FileSize:     tFile.FileSize.Int64,
	}
	return fMeta, nil
}

// deleteFileMeta 通过sha1值删除某个文件的元信息对象
func DeleteFileMeta(fileSha1 string) {
	delete(fileMetaMap, fileSha1)
}

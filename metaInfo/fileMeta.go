package metaInfo

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1, FileName, FileLocation, UpdateTime string
	FileSize                                     int64
}

var fileMetaMap map[string]FileMeta

func initFileMetaMap() {
	fileMetaMap = make(map[string]FileMeta)
}

// UpdateFileMetaMap 新增/更新文件元信息
func UpdateFileMetaMap(newMeta FileMeta) {
	fileMetaMap[newMeta.FileSha1] = newMeta
}

// GetFileMeta 通过sha1值获取某个文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetaMap[fileSha1]
}

// deleteFileMeta 通过sha1值删除某个文件的元信息对象
func deleteFileMeta(fileSha1 string) {
	delete(fileMetaMap, fileSha1)
}

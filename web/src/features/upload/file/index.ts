import { uploadFileToS3 } from './s3'

export async function uploadFile(file: File) {
  return uploadFileToS3(file)
}

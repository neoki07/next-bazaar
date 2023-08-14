import { S3 } from 'aws-sdk'

const s3 = new S3({
  accessKeyId: process.env.NEXT_PUBLIC_S3_ACCESS_KEY_ID,
  secretAccessKey: process.env.NEXT_PUBLIC_S3_SECRET_ACCESS_KEY,
  region: process.env.NEXT_PUBLIC_S3_REGION,
})

export async function uploadFileToS3(file: File) {
  if (process.env.NEXT_PUBLIC_S3_BUCKET_NAME === undefined) {
    throw new Error('S3_BUCKET_NAME is not defined')
  }

  const fileName = `${Date.now()}-${file.name}`
  const params: S3.PutObjectRequest = {
    Bucket: process.env.NEXT_PUBLIC_S3_BUCKET_NAME,
    Key: fileName,
    ContentType: file.type,
    Body: file,
  }

  const data = await s3.upload(params).promise()
  return data.Location
}

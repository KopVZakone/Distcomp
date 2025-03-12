package bsuir.dc.rest.service

import bsuir.dc.rest.dto.from.WriterRequestTo
import bsuir.dc.rest.dto.to.WriterResponseTo
import bsuir.dc.rest.mapper.toEntity
import bsuir.dc.rest.mapper.toResponse
import bsuir.dc.rest.repository.memory.IssueInMemoryRepository
import bsuir.dc.rest.repository.memory.WriterInMemoryRepository
import org.springframework.stereotype.Service

@Service
class WriterService(
    private val writerRepository: WriterInMemoryRepository,
    private val issueRepository: IssueInMemoryRepository,
) {
    fun createWriter(writerRequestTo: WriterRequestTo): WriterResponseTo {
        val writer = writerRequestTo.toEntity()
        val savedWriter = writerRepository.save(writer)
        return savedWriter.toResponse()
    }

    fun getWriterById(id: Long): WriterResponseTo {
        val writer = writerRepository.findById(id)
        return writer.toResponse()
    }

    fun getAllWriters(): List<WriterResponseTo> =
        writerRepository.findAll().map { it.toResponse() }

    fun updateWriter(id: Long, writerRequestTo: WriterRequestTo): WriterResponseTo {
        val updatedWriter = writerRequestTo.toEntity().apply { this.id = id }
        return writerRepository.update(updatedWriter).toResponse()
    }

    fun deleteWriter(id: Long) {
        writerRepository.deleteById(id)
    }

    fun getWriterByIssueId(issueId: Long): WriterResponseTo {
        val issue = issueRepository.findById(issueId)
        val writer = writerRepository.findById(issue.writerId)
        return writer.toResponse()
    }
}